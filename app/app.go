package app

import (
	"database/sql"
	"kaban-board-plus/common/component"
	"kaban-board-plus/screen/board"
	"kaban-board-plus/screen/list"

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

func (k KabanBoardPlus) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
            return k, nil
        case key.Matches(m, keys.BoardList):
            k.model = boardList
            k.screen = boardListScreen
            return k, nil
        }
    }
    var cmd tea.Cmd
    var s component.Screen
    if(k.model != nil) {
        s, cmd = k.model.Update(msg)
        k.model = s
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
