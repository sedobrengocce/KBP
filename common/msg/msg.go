package msg

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg struct {
    Err error
}

func NewErrorMsg(err error) tea.Cmd {
    return func() tea.Msg {
        return ErrorMsg{Err: err}
    }
}

type NewTaskMsg struct {
    Title string
    Description string
    priority int
}

func NewNewTaskMsg(title string, description string, priority int) tea.Cmd {
    return func() tea.Msg {
        return NewTaskMsg{Title: title, Description: description, priority: priority}
    }
}

type NewBoardMsg struct {
    Title string
}

func NewNewBoardMsg(title string) tea.Cmd {
    return func() tea.Msg {
        return NewBoardMsg{Title: title}
    }
}

type SelectBoardMsg struct {
    Id int
}

func NewSelectBoardMsg(id int) tea.Cmd {
    return func() tea.Msg {
        return SelectBoardMsg{Id: id}
    }
}
