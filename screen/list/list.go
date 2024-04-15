package list

import (
	"database/sql"
	"kaban-board-plus/common/component"
	mesg "kaban-board-plus/common/msg"
	"kaban-board-plus/component/button"
	"kaban-board-plus/component/dialog"
	dialogComponent "kaban-board-plus/component/dialog/components"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
    heightMargin = 6
    widthGrow = 3
)
var (
    BoardListColumnStyle = lipgloss.NewStyle().Padding(1,2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FF9800"))
)

type List struct {
    help help.Model
    showHelp bool
    db *sql.DB
    list list.Model
    preview list.Model
    height int
    width int
}

var dlg *dialog.Dialog

func (l *List) ToggleHelp() {
    l.showHelp = !l.showHelp
}

func NewList(db *sql.DB) *List {
    help := help.New()
    help.ShowAll = false
    l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Boards"
    l.SetShowHelp(false)
    preview := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
    preview.Title = ""
    preview.SetShowHelp(false)
    return &List{
        help: help,
        showHelp: false,
        db: db,
        list: l,
        preview: preview,
    }
}

func (l *List) getBoardList() tea.Cmd {
    rows, err := l.db.Query("SELECT id, name FROM boards")
    if err != nil {
        log.Print("Error getting boards: ", err)
        return mesg.NewErrorMsg(err)
    }
    defer rows.Close()
    var bl []boardElement
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        if err != nil {
            log.Print("Error getting boards: ", err)
        return mesg.NewErrorMsg(err)
        }
        b := newBoardElement(id, name)
        bl = append(bl, *b)
    }
    var newItems []list.Item
    for _, item := range bl {
        newItems = append(newItems, item)
    }
    cmd := l.list.SetItems(newItems)
    return cmd
}

func (l *List) deleteBoard(board_id int) tea.Cmd {
    _, err := l.db.Exec("DELETE FROM tasks WHERE board_id=?", board_id)
    if err != nil {
        log.Print("Error deleting tasks: ", err)
        return mesg.NewErrorMsg(err)
    }
    _, err = l.db.Exec("DELETE FROM boards WHERE id=?", board_id)
    if err != nil {
        log.Print("Error dele board: ", err)
        return mesg.NewErrorMsg(err)
    }
    return NewUpdateListMsg()
}

func (l *List) SetSize(width, height int) {
    l.width = width
    l.height = height - heightMargin
    l.list.SetSize(l.width, l.height)
}

func (l *List) Init() tea.Cmd {
    cmd := l.getBoardList()
    return cmd
}

func (l List) askDeleteBoard() *dialog.Dialog {
    item := l.list.SelectedItem()
    if item == nil {
        return nil
    }
    selectedItem := item.(boardElement)
    yesButton := button.NewButton("Yes", func() (tea.Cmd, error) {
        closeDialog()
        return NewDeleteBoardMsg(selectedItem.id), nil
    })
    noButton := button.NewButton("No", func() (tea.Cmd, error) {
        closeDialog()
        return nil, nil
    })
    message := dialogComponent.NewMessage("Are you sure you want to delete " + selectedItem.Title() + " board?")
    compnents := []dialog.DialogComponent{}
    compnents = append(compnents, message)
    noButton.Focus()
    d := dialog.NewDialog("Delete Board", compnents, 40, 10, 1, []button.Button{
        *yesButton,
        *noButton,
    })
    return d
}

func (l List) Update(msg tea.Msg) (component.Screen, tea.Cmd) {
    if dlg != nil {
        cmd := dlg.Update(msg)
        return &l, cmd
    }
    switch msg := msg.(type) {
        case UpdateListMsg:
            return &l, l.getBoardList()
        case DeleteBoardMsg:
            return &l, l.deleteBoard(msg.board_id)
        case tea.KeyMsg:
            switch {
                case key.Matches(msg, keys.Enter):
                    item := l.list.SelectedItem()
                    if item == nil {
                        return &l, nil
                    }
                    return &l, mesg.NewSelectBoardMsg(item.(boardElement).id)
                case key.Matches(msg, keys.DeleteBoard):
                    dlg = l.askDeleteBoard()
                    return &l, nil
                case key.Matches(msg, keys.NewBoard):
                    return &l, nil
                }
            }
    bl, cmd := l.list.Update(msg)
    l.list = bl
    return &l, cmd
}

func (l List) View() string {
    brdList := BoardListColumnStyle.
        Height(l.height).
        Width(l.width / widthGrow).
        Render(l.list.View())

    var help string
    if l.showHelp {
        help = l.help.View(keys)
    }

    // TODO: Implement preview
    preview := ""

    main := lipgloss.JoinHorizontal(
        lipgloss.Left,
        brdList,
        preview,
    )

    if dlg != nil {
        return lipgloss.Place(
            l.width,
            l.height,
            lipgloss.Center,
            lipgloss.Center,
            dlg.Render(),
        )
    }

    return lipgloss.JoinVertical(
        lipgloss.Left, 
        main,
        help,
    )
}

func closeDialog() {
    dlg = nil
}

