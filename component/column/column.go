package column

import (
	"kaban-board-plus/component/task"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    titleBorder = lipgloss.Border{
        Left: "",
        Right: "",
    }
    paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
    helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

var (
    focusedColumnTitleStyle = lipgloss.NewStyle().
        MarginLeft(2).
        Background(lipgloss.Color("#3f3f3f")).
        Foreground(lipgloss.Color("#ff9800")).
        Border(titleBorder, false, true).
        BorderForeground(lipgloss.Color("#3f3f3f")).
        Bold(true).
        Align(lipgloss.Center)
    focusedColumnBorderStyle = lipgloss.NewStyle().Padding(1,2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FF9800"))
)

var (
    unfocusedColumnTitleStyle = lipgloss.NewStyle().
        MarginLeft(2).
        Background(lipgloss.Color("#1a1a1a")).
        Foreground(lipgloss.Color("#3f3f3f")).
        Border(titleBorder, false, true).
        BorderForeground(lipgloss.Color("#1a1a1a")).
        Align(lipgloss.Center)
    unfocusedColumnBorderStyle = lipgloss.NewStyle().Padding(1,2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3f3f3f"))
)

const (
    heightMargin = 6
    widthMargin = 2
)

type Column struct { 
    title string
    taskList list.Model
    focused bool
    width int
    height int
}

func NewColumn(title string) *Column {
    l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
    l.Title = " " + title + " "
    l.Styles.Title = unfocusedColumnTitleStyle
    l.SetShowHelp(false)
    return &Column{
        title: " " + title + " ",
        focused: false,
        taskList: l,
    }
}

func (c *Column) AddItems(items []task.Task) tea.Cmd {
    l := []list.Item{}
    for _, item := range items {
        l = append(l, item)
    }
    return c.taskList.SetItems(l)
}

func (c *Column) AddItem(item task.Task) tea.Cmd {
    return c.taskList.InsertItem(len(c.taskList.Items()), item)
}

func (c *Column) RemoveItem() {
    c.taskList.RemoveItem(c.taskList.Index())
}

func (c Column) Length() int {
    return len(c.taskList.Items())
}

func (c Column) SelectedItem() task.Task {
    return c.taskList.SelectedItem().(task.Task)
}

func (c *Column) Focus() {
    c.focused = true
    c.taskList.Styles.Title = focusedColumnTitleStyle
}

func (c *Column) Blur() {
    c.focused = false
    c.taskList.Styles.Title = unfocusedColumnTitleStyle
}

func (c Column) Focused() bool {
    return c.focused
}

func (c *Column) SetSizes(w, h int) {
    c.width = w - widthMargin
    c.height = h - heightMargin
    c.taskList.SetSize(w - widthMargin, h - heightMargin)
}

func (c *Column) Update(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    c.taskList, cmd = c.taskList.Update(msg)
    return cmd
}

func (c Column) View() string {
    if(c.focused) {
        return focusedColumnBorderStyle.
            Width(c.width).
            Height(c.height).
            Render(c.taskList.View())
    }
    return unfocusedColumnBorderStyle.
            Width(c.width).
            Height(c.height).
            Render(c.taskList.View())
}
