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
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
    showArchived bool
    focusedColumn int
    height int
    width int
    isToday bool
    cols []column.Column
    db   *sql.DB
}

type boardListElement struct {
    id int
    name string
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
        showArchived: false,
        isToday: false,
    }
}

func NewTodayBoard(db *sql.DB) *Board {
    b := newBoard(0, "Today", db)
    b.cols = append(b.cols, *column.NewColumn("To Do", true))
    b.cols = append(b.cols, *column.NewColumn("In Progress", true))
    b.cols = append(b.cols, *column.NewColumn("Done", true))
    b.isToday = true
    return b
}

func NewProjectBoard(id int, name string, db *sql.DB) *Board {
    b := newBoard(id, name, db)
    b.cols = append(b.cols, *column.NewColumn("To Do", false))
    b.cols = append(b.cols, *column.NewColumn("Done", false))
    b.showArchived = true
    return b
}

func (b *Board) getTasks(archived bool) ([]task.Task, error) {
    var tasks []task.Task
    var rows *sql.Rows
    var err error
    if b.isToday {
        rows, err = b.db.Query("SELECT id, name, description, priority, status, is_archived, is_today FROM tasks WHERE is_today = true and is_archived = false")
    } else {
        rows, err = b.db.Query("SELECT id, name, description, priority, status, is_archived, is_today FROM tasks WHERE board_id = ?", b.id)
    }
    defer rows.Close()
    if err != nil {
        log.Print("Error getting tasks: ", err)
        return tasks, err
    }
    for rows.Next() {
        var id int
        var title string
        var description string
        var priority int
        var status int
        var isArchived bool
        var isToday bool
        err = rows.Scan(&id, &title, &description, &priority, &status, &isArchived, &isToday)
        if err != nil {
            log.Print("Error scanning row: ", err)
            return tasks, err
        }
        if archived { 
            t := task.NewTask(id, title, description, task.Priority(priority), task.Status(status), isArchived, isToday)
            tasks = append(tasks, t)
        } else {
            if !isArchived {
                t := task.NewTask(id, title, description, task.Priority(priority), task.Status(status), isArchived, isToday)
                tasks = append(tasks, t)
            }
        }
    }
    return tasks, nil
}

func (b *Board) setTasks(tasks *[]task.Task) tea.Cmd {
    tasksDone := []task.Task{}
    tasksTodo := []task.Task{}
    tasksInProgress := []task.Task{}
    for _, t := range *tasks {
        switch t.Status() {
        case task.Todo:
            tasksTodo = append(tasksTodo, t)
        case task.InProgress:
            tasksInProgress = append(tasksInProgress, t)
        case task.Done:
            tasksDone = append(tasksDone, t)
        default:
        }
    }
    if b.isToday {
        return tea.Batch(
            b.cols[task.Todo].SetItems(tasksTodo),
            b.cols[task.InProgress].SetItems(tasksInProgress),
            b.cols[task.Done].SetItems(tasksDone),
        )
    } else {
        composedTodo := append(tasksTodo, tasksInProgress...)
        return tea.Batch(
            b.cols[task.Todo].SetItems(composedTodo),
            b.cols[task.InProgress].SetItems(tasksDone),
        )
    }
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
    return func() tea.Msg {
        if (b.cols[b.focusedColumn].Length() == 0) {
            return msgs.ErrorMsg{Err: NewEmptyColumnError()}
        }
        nextColumn := (b.focusedColumn + 1) % len(b.cols)
        selectedTask := b.cols[b.focusedColumn].SelectedItem()

        _, err := b.db.Exec("UPDATE tasks SET status = ? WHERE id = ?", nextColumn, selectedTask.ID())
        if err != nil {
            log.Print("Error updating task: ", err)
            return msgs.ErrorMsg{Err: NewDbError(err)}
        }
        if nextColumn == int(task.Done) {
            _, err := b.db.Exec("UPDATE tasks SET done_date = ? WHERE id = ?", time.Now().Local().Format("2006-01-02"), selectedTask.ID())
            if err != nil {
                log.Print("Error updating task: ", err)
                return msgs.ErrorMsg{Err: NewDbError(err)}
            }
        } else if nextColumn == int(task.Todo) {
            _, err := b.db.Exec("UPDATE tasks SET done_date = NULL WHERE id = ?", selectedTask.ID())
            if err != nil {
                log.Print("Error updating task: ", err)
                return msgs.ErrorMsg{Err: NewDbError(err)}
            }
        }

        return UpdateMsg{}
    }
}

func (b *Board) archiveAll() tea.Cmd {
    return func() tea.Msg {
        _, err := b.db.Exec("UPDATE tasks SET is_archived = true WHERE board_id = ? and status = 2", b.id)
        if err != nil {
            log.Print("Error archiving all tasks: ", err)
            return msgs.ErrorMsg{Err: NewDbError(err)}
        }
        return UpdateMsg{}
    }
}

func (b *Board) toggleArchive() tea.Cmd {
    return func() tea.Msg {
        if b.cols[b.focusedColumn].Length() == 0 {
            log.Print("Empty column")
            return msgs.ErrorMsg{Err: NewEmptyColumnError()}
        }
        selectedTask := b.cols[b.focusedColumn].SelectedItem()
        if selectedTask.Status() != task.Done {
            return nil
        }
        _, err := b.db.Exec("UPDATE tasks SET is_archived = ? WHERE id = ?", !selectedTask.IsArchived(), selectedTask.ID())
        if err != nil {
            log.Print("Error updating task: ", err)
            return msgs.ErrorMsg{Err: NewDbError(err)}
        }
        return UpdateMsg{}
    }
}

func (b *Board) toggleToday() tea.Cmd {
    return func() tea.Msg {
        if (b.cols[b.focusedColumn].Length() == 0) {
            log.Print("Empty column")
            return msgs.ErrorMsg{Err: NewEmptyColumnError()}
        }
        selectedTask := b.cols[b.focusedColumn].SelectedItem()
        _, err := b.db.Exec("UPDATE tasks SET is_today = ? WHERE id = ?", !selectedTask.IsToday(), selectedTask.ID())
        if err != nil {
            log.Print("Error updating task: ", err)
            return msgs.ErrorMsg{Err: NewDbError(err)}
        }
        return UpdateMsg{}
    }
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
    confirmButton := button.NewButton("Confirm", func() tea.Cmd {
        closeDialog()
        return NewCreateTaskMsg(title.GetText(), description.GetText(), radioInput.GetValue())
    })
    cancelButton := button.NewButton("Cancel", func() tea.Cmd {
        closeDialog()
        return nil
    })
    buttons := []button.Button{
        *confirmButton,
        *cancelButton,
    }
    d := dialog.NewDialog("New Task", components, 40, 10, buttons)
    dlg = d
    return nil
}

func (b *Board) askDeleteTask() tea.Cmd {
    if b.cols[b.focusedColumn].Length() == 0 {
        return nil
    }
    item := b.cols[b.focusedColumn].SelectedItem()
    yesButton := button.NewButton("Yes", func() tea.Cmd {
        return func() tea.Msg {
            closeDialog()
            return deleteMsg{id: item.ID()}
        }
    })
    noButton := button.NewButton("No", func() tea.Cmd {
        return func() tea.Msg {
            closeDialog()
            return nil
        }
    })
    message := dialogComponent.NewMessage("Are you sure you want to delete " + item.Title() + " task?")
    noButton.Focus()
    d := dialog.NewDialog("Delete Task", 
        []dialog.DialogComponent{message}, 
        40, 10,
        []button.Button{
            *yesButton,
            *noButton,
        },
    )
    dlg = d
    return nil
}

func (b* Board) askMoveTaskToBoard() tea.Cmd {
    if b.cols[b.focusedColumn].Length() == 0 {
        return nil
    }
    item := b.cols[b.focusedColumn].SelectedItem()
    message := dialogComponent.NewMessage("Please select destination Board")
    radioItems := []dialogComponent.RadioItem[int]{}
    list, err := b.getOtherBoards()
    if err != nil {
        return func() tea.Msg {
            return msgs.ErrorMsg{Err: err}
        }
    }
    for _,l := range list {
        i := dialogComponent.NewRadioItem(l.name, l.id)
        radioItems = append(radioItems, *i)
    }
    boardList := dialogComponent.NewRadioInput(
        "Boards: ",
        radioItems,
    )
    boardList.Focus()
    yesButton := button.NewButton("Yes", func() tea.Cmd {
        return func() tea.Msg {
            closeDialog()
            return moveTaskToBoardMsg{id: item.ID(), boardId: boardList.GetValue()}
        }
    })
    noButton := button.NewButton("No", func() tea.Cmd {
        return func() tea.Msg {
            closeDialog()
            return nil
        }
    })
    d := dialog.NewDialog("Move Task", 
        []dialog.DialogComponent{message, boardList}, 
        40, 10,
        []button.Button{
            *yesButton,
            *noButton,
        },
    )
    dlg = d
    return nil
}

func (b Board) getOtherBoards() ([]boardListElement, error) {
    boardList := []boardListElement{}

    var rows *sql.Rows
    var err error
    rows, err = b.db.Query("SELECT id, name FROM boards WHERE id != ?", b.id)

    if err != nil {
        log.Print("Error getting boards: ", err)
        return boardList, err
    }
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        if err != nil {
            log.Print("Error scanning row: ", err)
            return boardList, err
        }
        element := boardListElement{id: id, name:name}
        boardList = append(boardList, element)
    }
    return boardList, nil
}

func (b Board) deleteTask(id int) tea.Cmd {
    return func() tea.Msg {
        _, err := b.db.Exec("DELETE from tasks WHERE id = ?", id)
        if err != nil {
            log.Print("Error dedleting task: ", id)
            return dbError{err: err}
        }
        return UpdateMsg{}
    }
}

func (b Board) moveTaskToBoard(taskid, boardid int) tea.Cmd {
    _, err := b.db.Exec("UPDATE tasks SET board_id = ? where id = ?", boardid, taskid)
    if err != nil {
        return func() tea.Msg {
            log.Print("Cannot chage board to task due error: ", err)
            return msgs.ErrorMsg{Err: err}
        }
    }
    return func() tea.Msg {
        return UpdateMsg{}
    }
}

func (b Board) createTask(name, description string, priority task.Priority) tea.Cmd {
    _, err := b.db.Exec("INSERT INTO tasks (board_id, name, description, priority) VALUES (?,?,?,?)", b.id, name, description, int(priority))
    if err != nil {
        log.Print("Error creating new task: ", err)
        return msgs.NewErrorMsg(err)
    }
    return NewUpdateMsg()
}

func (b *Board) showHideArchive() tea.Cmd {
    b.showArchived = !b.showArchived
    tasks, err := b.getTasks(b.showArchived)     
    if err != nil {
        return func() tea.Msg {
            return msgs.ErrorMsg{Err: err}
        }
    }
    return b.setTasks(&tasks)
}

func (b *Board) Init() tea.Cmd {
    tasks, err := b.getTasks(b.showArchived)
    if err != nil {
        return func() tea.Msg {
            return msgs.ErrorMsg{Err: err}
        }
    }
    if b.focusedColumn == 0 {
        b.cols[0].Focus()
    }
    return b.setTasks(&tasks)
}

func (b *Board) Update(msg tea.Msg) (component.Screen, tea.Cmd) {
    if dlg != nil {
        return b, dlg.Update(msg)
    }
    switch msg := msg.(type) {
    case UpdateMsg:
        return b, b.Init() 
    case CreateTaskMsg:
        return b, b.createTask(msg.Name, msg.Description, msg.Priority)
    case deleteMsg:
        return b, b.deleteTask(msg.id)
    case moveTaskToBoardMsg:
        return b, b.moveTaskToBoard(msg.id, msg.boardId)
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Left):
            b.prevColumn()
        case key.Matches(msg, keys.Right):
            b.nextColumn()
        case key.Matches(msg, keys.Action):
            if b.isToday {
                return b, b.moveTask()
            } 
            return b, b.toggleToday()
        case key.Matches(msg, keys.NewTask):
            if !b.isToday {
                return b, b.askNewTask()
            }            
        case key.Matches(msg, keys.Archive):
            return b, b.toggleArchive()
        case key.Matches(msg, keys.ShowHideArchive):
            if !b.isToday {
                return b, b.showHideArchive()
            }
        case key.Matches(msg, keys.ArchiveAll):
            return b, b.archiveAll()
        case key.Matches(msg, keys.Delete):
            return b, b.askDeleteTask()
        case key.Matches(msg, keys.Move):
            if !b.isToday {
                return b, b.askMoveTaskToBoard()
            }            
        }
    }

    b.cols[b.focusedColumn].Update(msg)

    return b, nil
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
