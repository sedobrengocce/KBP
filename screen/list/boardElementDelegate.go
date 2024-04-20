package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff8f50"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff9800")).Bold(true).
                            BorderStyle(lipgloss.ThickBorder()).BorderLeft(true).BorderForeground(lipgloss.Color("#ff9800"))
)

type boardElementDelegate struct {} 

func (d boardElementDelegate) Height() int                             { return 1 }
func (d boardElementDelegate) Spacing() int                            { return 0 }
func (d boardElementDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d boardElementDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(boardElement)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
        fn = selectedItemStyle.Render
	}

	fmt.Fprint(w, fn(i.title()))
}
