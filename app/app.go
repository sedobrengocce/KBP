// TODO: Refactor Help for entire app
package app

import (
	"database/sql"
	"kaban-board-plus/common/component"
	mesg "kaban-board-plus/common/msg"
	"kaban-board-plus/screen/board"
	"kaban-board-plus/screen/list"
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
    todayScreen screen = iota
    boardListScreen
    boardScreen
)

var todayBoard *board.Board
var boardList *list.List

type initialMsg struct{}

func newInitialMsg() tea.Cmd {
    return func() tea.Msg {
        return initialMsg{}
    }
}

type KabanBoardPlus struct {
    screen screen
    model  component.Screen
    height int
    width  int
    help   help.Model
    db     *sql.DB
    quitting bool
    err     error
}

func NewKabanBoardPlus(db *sql.DB) *KabanBoardPlus {
    help := help.New()
    help.ShowAll = false
    todayBoard = board.NewTodayBoard(db)
    return &KabanBoardPlus{
        db: db,
        help: help,
        quitting: false,
        err: nil,
    }
}

func (k *KabanBoardPlus) SetSize(width, height int) {
    k.width = width
    k.height = height
}

func (k KabanBoardPlus) Init() tea.Cmd {
    return newInitialMsg()
}

func (k KabanBoardPlus) getBoardName(id int) (string, error) {
    rows, err := k.db.Query("SELECT name FROM boards WHERE id=?", id)
    if err != nil {
        log.Print("Error retrieving board name: ", err)
        return "", err
    }
     
    defer rows.Close()

    var name string
    for rows.Next() {
        err = rows.Scan(&name)
        if err != nil {
            log.Print("Error reading board name: ", err)
            return "", err
        }
    }

    return name, nil
}

func (k *KabanBoardPlus) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var s component.Screen
    if(k.model != nil) {
        s, cmd = k.model.Update(msg)
        k.model = s
        if cmd != nil {
            return k, cmd
        }
    }
    switch m := msg.(type) {
    case initialMsg:
        var cmd tea.Cmd
        var cmds []tea.Cmd
        todayBoard = board.NewTodayBoard(k.db)
        cmd = todayBoard.Init()
        cmds = append(cmds, cmd)
        boardList = list.NewList(k.db)
        cmd = boardList.Init()
        cmds = append(cmds, cmd)
        if(k.height > 0 && k.width > 0) {
            todayBoard.SetSize(k.width, k.height)
            boardList.SetSize(k.width, k.height)
        }
        k.screen = todayScreen
        k.model = todayBoard 
        return k, tea.Batch(cmds...)
    case mesg.SelectBoardMsg:
        name, err := k.getBoardName(m.Id)
        if err != nil {
            return k, mesg.NewErrorMsg(err)
        }
        brd := board.NewProjectBoard(m.Id, name, k.db)
        cmd := brd.Init()
        brd.SetSize(k.width, k.height)
        k.screen = boardListScreen
        k.model = brd
        return k, cmd
    case tea.WindowSizeMsg:
        if(todayBoard != nil) {
            todayBoard.SetSize(m.Width, m.Height)
        }
        if (boardList != nil) {
            boardList.SetSize(m.Width, m.Height)
        }
        k.SetSize(m.Width, m.Height)
        return k, nil
    case tea.KeyMsg:
        switch {
        case key.Matches(m, keys.Quit):
            k.quitting = true;
            return k, tea.Quit
        case key.Matches(m, keys.Help):
            k.model.ToggleHelp()
            return k, nil
        case key.Matches(m, keys.Today):
            k.model = todayBoard
            k.screen = todayScreen
            k.model.Init()
            return k, nil
        case key.Matches(m, keys.BoardList):
            k.model = boardList
            k.screen = boardListScreen
            k.model.Init()
            return k, nil
        }
    }
    return k, cmd
}

func (k KabanBoardPlus) View() string {
    if(k.quitting) {
        return "See you next time!"
    } else if (k.model == nil) {
        return "Loading..."
    }
    return lipgloss.JoinVertical(
        lipgloss.Left,
        k.model.View(),
        k.help.View(keys),
        //lipgloss.PlaceVertical(1, lipgloss.Bottom, k.help.View(keys)),
    )
}
