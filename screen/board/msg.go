package board

import (
	"kaban-board-plus/component/task"

	tea "github.com/charmbracelet/bubbletea"
)

type UpdateMsg struct {}

func NewUpdateMsg() tea.Cmd {
    return func () tea.Msg {
        return UpdateMsg{}
    }
}

type CreateTaskMsg struct {
    Name string
    Description string
    Priority task.Priority
}

func NewCreateTaskMsg(name, description string, priority task.Priority) tea.Cmd {
    return func() tea.Msg {
        return CreateTaskMsg{
            Name: name,
            Description: description,
            Priority: priority,
        }
    }
}

type SetTodoTaskMsg struct {
    tasks   *[]task.Task
}

type SetProgressTaskMsg struct {
    tasks   *[]task.Task
}

type SetDoneTaskMsg struct {
    tasks   *[]task.Task
}

type deleteMsg struct {
    id int
}
