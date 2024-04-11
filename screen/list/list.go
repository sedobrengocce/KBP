package list

import (
	"database/sql"
	"kaban-board-plus/common/component"
	mesg "kaban-board-plus/common/msg"
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
    var rows *sql.Rows
    var err error
    rows, err = l.db.Query("SELECT id, name FROM boards")
    if err != nil {
        log.Print("Error getting boards: ", err)
        return mesg.NewErrorMsg(err)
    }
    defer rows.Close()
    var bl []list.Item
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        log.Print("Board: ", id, name)
        if err != nil {
            log.Print("Error getting boards: ", err)
        return mesg.NewErrorMsg(err)
        }
        b := newBoardElement(id, name)
        bl = append(bl, b)
    }
    cmd := l.list.SetItems(bl)
    return cmd
}

func (l *List) SetSize(width, height int) {
    l.width = width / widthGrow
    l.height = height - heightMargin
    l.list.SetSize(l.width, l.height)
}

func (l *List) Init() tea.Cmd {
    cmd := l.getBoardList()
    return cmd
}

func (l List) Update(msg tea.Msg) (component.Screen, tea.Cmd) {
    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
                case key.Matches(msg, keys.Enter):
                    item := l.list.SelectedItem()
                    if item == nil {
                        return &l, nil
                    }
                    return &l, mesg.NewSelectBoardMsg(item.(boardElement).id)
                case key.Matches(msg, keys.DeleteBoard):
                    return &l, nil
                case key.Matches(msg, keys.NewBoard):
                    return &l, nil
                }
            }
    _, cmd := l.list.Update(msg)
    return &l, cmd
}

func (l List) View() string {
    brdList := BoardListColumnStyle.
        Height(l.height).
        Width(l.width).
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

    return lipgloss.JoinVertical(
        lipgloss.Left, 
        main,
        help,
    )
}

