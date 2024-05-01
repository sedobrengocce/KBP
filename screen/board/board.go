// TODO: Move task to another board
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
    showArchived bool
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
        showArchived: false,
    }
}

func NewTodayBoard(db *sql.DB) *Board {
    b := newBoard(0, "Today", db)
    b.cols = append(b.cols, *column.NewColumn("To Do", true))
    b.cols = append(b.cols, *column.NewColumn("In Progress", true))
    b.cols = append(b.cols, *column.NewColumn("Done", true))
    return b
}

func NewProjectBoard(id int, name string, db *sql.DB) *Board {
    b := newBoard(id, name, db)
    b.cols = append(b.cols, *column.NewColumn("To Do", false))
    b.cols = append(b.cols, *column.NewColumn("Done", false))
    b.showArchived = true
    return b
}

func (b *Board) getTasks(status task.Status, today, archived bool) tea.Cmd {
    return func() tea.Msg {
        if status == task.InProgress && !today {
            return nil
        }
        var tasks []task.Task
        var rows *sql.Rows
        var err error
        if today {
            rows, err = b.db.Query("SELECT id, name, description, priority, is_archived, is_today FROM tasks WHERE status = ? and is_today = true and is_archived = false", status)
        } else {
            if status == task.Done {
                rows, err = b.db.Query("SELECT id, name, description, priority, is_archived, is_today FROM tasks WHERE status = 2 and board_id = ?", b.id)
            } else {
                rows, err = b.db.Query("SELECT id, name, description, priority, is_archived, is_today FROM tasks WHERE (status = 0 or status = 1) and board_id = ?", b.id)
            }
        }
        defer rows.Close()
        if err != nil {
            log.Print("Error getting tasks: ", err)
            return msgs.ErrorMsg{Err: err}
        }
        for rows.Next() {
            var id int
            var title string
            var description string
            var priority int
            var isArchived bool
            var isToday bool
            err = rows.Scan(&id, &title, &description, &priority, &isArchived, &isToday)
            if err != nil {
                log.Print("Error scanning row: ", err)
                return msgs.ErrorMsg{Err: err}
            }
            if archived { 
                t := task.NewTask(id, title, description, task.Priority(priority), status, isArchived, isToday)
                tasks = append(tasks, t)
            } else {
                if !isArchived {
                    t := task.NewTask(id, title, description, task.Priority(priority), status, isArchived, isToday)
                    tasks = append(tasks, t)
                }
            }
        }
        switch status {
        case task.Todo:
            return SetTodoTaskMsg{tasks: &tasks}
        case task.InProgress:
            return SetProgressTaskMsg{tasks: &tasks}
        case task.Done:
            if today {
                return SetDoneTaskMsg{tasks: &tasks}
            }  
            return SetProgressTaskMsg{tasks: &tasks}
        default:
            log.Print("Unknow status type: ", status)
            return msgs.ErrorMsg{Err: newStatusTypeError(status)}
        }
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
    return b.getTasks(task.Done, false, b.showArchived)     
}

func (b *Board) Init() tea.Cmd {
    var cmds = []tea.Cmd{}
    for i := 0; i < int(task.Unknown) ; i++ {
        cmd := b.getTasks(task.Status(i), len(b.cols) == todayColumnNum, b.showArchived)
        cmds = append(cmds, cmd)
    } 
    if b.focusedColumn == 0 {
        b.cols[0].Focus()
    }
    return tea.Batch(cmds...)
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
    case SetDoneTaskMsg:
        return b, b.cols[task.Done].SetItems(*msg.tasks)
    case SetProgressTaskMsg:
        return b, b.cols[task.InProgress].SetItems(*msg.tasks)
    case SetTodoTaskMsg:
        return b, b.cols[task.Todo].SetItems(*msg.tasks)
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Left):
            b.prevColumn()
        case key.Matches(msg, keys.Right):
            b.nextColumn()
        case key.Matches(msg, keys.Action):
            if(len(b.cols) == todayColumnNum) {
                return b, b.moveTask()
            } 
            return b, b.toggleToday()
        case key.Matches(msg, keys.NewTask):
            if(len(b.cols) == projectColumnNum) {
                return b, b.askNewTask()
            }            
        case key.Matches(msg, keys.Archive):
            return b, b.toggleArchive()
        case key.Matches(msg, keys.ShowHideArchive):
            if len(b.cols) == projectColumnNum {
                return b, b.showHideArchive()
            }
        case key.Matches(msg, keys.ArchiveAll):
            return b, b.archiveAll()
        case key.Matches(msg, keys.Delete):
            return b, b.askDeleteTask()
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
