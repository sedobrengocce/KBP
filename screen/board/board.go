package board

import (
	"database/sql"
	"kaban-board-plus/common/component"
	msgs "kaban-board-plus/common/msg"
	"kaban-board-plus/component/button"
	"kaban-board-plus/component/column"
	"kaban-board-plus/component/dialog"
	dialogComponent "kaban-board-plus/component/dialog/components"
	"kaban-board-plus/component/task"
	"log"
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
    todayColumnNum = 3
    projectColumnNum = 2
)

const (
    topMargin = 2
)

var (
    titleStyle = lipgloss.NewStyle().Align(lipgloss.Center).Background(lipgloss.Color("#3F3F3F")).Foreground(lipgloss.Color("#FF9800"))
    borderTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3F3F3F"))
)

var dlg *dialog.Dialog

type Board struct {
    id   int
    name string
    help help.Model
    showHelp bool
    focusedColumn int
    height int
    width int
    cols []column.Column
    db   *sql.DB
}

func newBoard(id int, name string, db *sql.DB) *Board {
    help := help.New()
    help.ShowAll = false
    return &Board{
        id: id,
        help: help,
        showHelp: false,
        name: name,
        focusedColumn: 0,
        db: db,
    }
}

func NewTodayBoard(db *sql.DB) *Board {
    b := newBoard(0, "Today", db)
    b.cols = append(b.cols, *column.NewColumn("To Do"))
    b.cols = append(b.cols, *column.NewColumn("In Progress"))
    b.cols = append(b.cols, *column.NewColumn("Done"))
    return b
}

func NewProjectBoard(id int, name string, db *sql.DB) *Board {
    b := newBoard(id, name, db)
    b.cols = append(b.cols, *column.NewColumn("To Do"))
    b.cols = append(b.cols, *column.NewColumn("Done"))
    return b
}

func (b *Board) getTasks(status task.Status, today bool) ([]task.Task, error) {
    var tasks []task.Task
    var rows *sql.Rows
    var err error
    if today {
        rows, err = b.db.Query("SELECT id, name, description, priority, is_archived FROM tasks WHERE status = ? and is_today = true and is_archived = false", status)
    } else {
        rows, err = b.db.Query("SELECT id, name, description, priority, is_archived FROM tasks WHERE status = ? and board_id = ?", status, b.id)
    }
    if err != nil {
        log.Print("Error getting tasks: ", err)
        return nil, err
    }
    for rows.Next() {
        var id int
        var title string
        var description string
        var priority int
        var isArchived bool
        err = rows.Scan(&id, &title, &description, &priority, &isArchived)
        if err != nil {
            log.Print("Error scanning row: ", err)
            return nil, err
        }
        t := task.NewTask(id, title, description, task.Priority(priority), status, isArchived)
        tasks = append(tasks, t)
    }
    defer rows.Close()
    return tasks, nil
}

func (b *Board) SetSize(width, height int) {
    b.width = width
    b.height = height - topMargin
    for i := range b.cols {
        b.cols[i].SetSizes(b.width/len(b.cols), b.height)
    }
}

func (b *Board) ToggleHelp() {
    b.showHelp = !b.showHelp
}

func (b *Board) prevColumn() {
    b.cols[b.focusedColumn].Blur()
    b.focusedColumn = (b.focusedColumn - 1 + len(b.cols)) % len(b.cols)
    b.cols[b.focusedColumn].Focus()
}

func (b *Board) nextColumn() {
    b.cols[b.focusedColumn].Blur()
    b.focusedColumn = (b.focusedColumn + 1) % len(b.cols) 
    b.cols[b.focusedColumn].Focus()
}

func (b *Board) moveTask() tea.Cmd {
    if (b.cols[b.focusedColumn].Length() == 0) {
        return msgs.NewErrorMsg(&emptyColumn{})
    }
    nextColumn := (b.focusedColumn + 1) % len(b.cols)
    selectedTask := b.cols[b.focusedColumn].SelectedItem()

    _, err := b.db.Exec("UPDATE tasks SET status = ? WHERE id = ?", nextColumn, selectedTask.ID())
    if err != nil {
        log.Print("Error updating task: ", err)
        return msgs.NewErrorMsg(&dbError{err: err})
    }

    return NewUpdateMsg()
}

func (b *Board) askNewTask() tea.Cmd {
    title := dialogComponent.NewTextInput("Task Name: ", "Name")
    title.Focus()
    description := dialogComponent.NewTextArea("Description: ", "Description")
    prioHigh := dialogComponent.NewRadioItem("High", task.High)
    prioMedium := dialogComponent.NewRadioItem("Medium", task.Medium)
    prioLow := dialogComponent.NewRadioItem("Low", task.Low)
    radioInput := dialogComponent.NewRadioInput(
        "Priority: ",
        []dialogComponent.RadioItem[task.Priority]{*prioHigh, *prioMedium, *prioLow},
    )
    components := []dialog.DialogComponent{
        title,
        description,
        radioInput,
    }
    confirmButton := button.NewButton("Confirm", func() (tea.Cmd, error) {
        closeDialog()
        return NewCreateTaskMsg(title.GetText(), description.GetText(), radioInput.GetValue()), nil
    })
    cancelButton := button.NewButton("Cancel", func() (tea.Cmd, error) {
        closeDialog()
        return nil, nil
    })
    buttons := []button.Button{
        *confirmButton,
        *cancelButton,
    }
    d := dialog.NewDialog("New Task", components, 40, 10, buttons)
    dlg = d
    return nil
}

func (b Board) createTask(name, description string, priority task.Priority) error {
    _, err := b.db.Exec("INSERT INTO tasks (board_id, name, description, priority) VALUES (?,?,?,?)", b.id, name, description, int(priority))
    return err
}

func (b *Board) Init() tea.Cmd {
    colNum := len(b.cols)
    var tasksList [][]task.Task
    if colNum == todayColumnNum {
        for i := 0; i < colNum; i++ {
            tasks, err := b.getTasks(task.Status(i), true)
            if err != nil {
                return  msgs.NewErrorMsg(err)
            }
            tasksList = append(tasksList, tasks)
        } 
    } else if colNum == projectColumnNum {
        var todos []task.Task
        var inProgress []task.Task
        var done []task.Task
        var err error
        todos, err = b.getTasks(task.Todo, false)
        if err != nil {
            return msgs.NewErrorMsg(err)
        }
        inProgress, err = b.getTasks(task.InProgress, false)
        if err != nil {
            return  msgs.NewErrorMsg(err)
        }
        done, err = b.getTasks(task.Done, false)
        if err != nil {
            return msgs.NewErrorMsg(err)
        }
        tasksList = append(tasksList, slices.Concat(todos, inProgress))
        tasksList = append(tasksList, done)

    } else {
        return msgs.NewErrorMsg(&errorColNumber{number: colNum})
    }
    var cmds []tea.Cmd
    for i := range b.cols {
        cmd := b.cols[i].AddItems(tasksList[i])
        cmds = append(cmds, cmd)
    }
    b.cols[b.focusedColumn].Focus()
    return tea.Batch(cmds...)
}

func (b Board) Update(msg tea.Msg) (component.Screen, tea.Cmd) {
    if dlg != nil {
        return &b, dlg.Update(msg)
    }
    switch msg := msg.(type) {
    case UpdateMsg:
        return &b, b.Init() 
    case CreateTask:
        err := b.createTask(msg.Name, msg.Description, msg.Priority)
        if err != nil {
            return &b, msgs.NewErrorMsg(err)
        }
        return &b, NewUpdateMsg()
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Left):
            b.prevColumn()
        case key.Matches(msg, keys.Right):
            b.nextColumn()
        case key.Matches(msg, keys.Enter):
            if(len(b.cols) == todayColumnNum) {
                return &b, b.moveTask()
            }
        case key.Matches(msg, keys.NewTask):
            if(len(b.cols) == projectColumnNum) {
                return &b, b.askNewTask()
            }            
        }
    }

    b.cols[b.focusedColumn].Update(msg)

    return &b, nil
}

func (b Board) View() string {
    if dlg != nil {
        return lipgloss.Place(
            b.width,
            b.height,
            lipgloss.Center,
            lipgloss.Center,
            dlg.Render(),
        )
    }
    var renderdColumns []string
    for i := range b.cols {
        renderdColumns = append(renderdColumns, b.cols[i].View())
    }
    titleLeftBorder := borderTitleStyle.Render("")
    titleRightBorder := borderTitleStyle.Render("")
    title := lipgloss.JoinHorizontal(lipgloss.Center, titleLeftBorder, titleStyle.Render(" " + b.name + " "), titleRightBorder)
    titleLine := lipgloss.PlaceHorizontal(b.width, titleStyle.GetAlign(), title)
    brd := lipgloss.JoinHorizontal(lipgloss.Top, renderdColumns...)
    h := ""
    if (b.showHelp) {
        h = b.help.View(keys)
    }
    return lipgloss.JoinVertical(lipgloss.Top, titleLine, brd, h)
}

func closeDialog() {
    dlg = nil
}
