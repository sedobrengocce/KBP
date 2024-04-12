package list

import tea "github.com/charmbracelet/bubbletea"

type UpdateListMsg struct {}

func NewUpdateListMsg() tea.Cmd {
    return func() tea.Msg {
        return UpdateListMsg{}
    }
}

type DeleteBoardMsg struct {
    board_id int
}

func NewDeleteBoardMsg(id int) tea.Cmd {
    return func() tea.Msg {
        return DeleteBoardMsg{board_id: id}
    }
}

type CreateBoardMsg struct {
    title string
}

func NewCreateBoardMsg(title string) tea.Cmd {
    return func() tea.Msg {
        return CreateBoardMsg{title: title}
    }
}
