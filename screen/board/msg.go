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

type CreateTask struct {
    Name string
    Description string
    Priority task.Priority
}

func NewCreateTaskMsg(name, description string, priority task.Priority) tea.Cmd {
    return func() tea.Msg {
        return NewCreateTaskMsg(name, description, priority)
    }
}
