package column

import (
	"fmt"
	"io"
	"kaban-board-plus/component/task"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    unselectedPriorityStyle = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).Width(3).Height(3).
                        Foreground(lipgloss.Color("#ff8550")).
                        BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff8550"))
    unselectedTitleStyle = lipgloss.NewStyle().Height(2).
                        Foreground(lipgloss.Color("#ff8550"))
    unselectedDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8f50"))
)

var (
    selectedPriorityStyle = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).Width(3).Height(3).
                        Bold(true).Foreground(lipgloss.Color("#ff9800")).Background(lipgloss.Color("#3f3f3f")).
                        BorderStyle(lipgloss.ThickBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff9800")).BorderBackground(lipgloss.Color("#3f3f3f"))
    selectedTitleStyle = lipgloss.NewStyle().Height(2).
                        Bold(true).Foreground(lipgloss.Color("#ff9800")).Background(lipgloss.Color("#3f3f3f"))
    selectedDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8f50")).Background(lipgloss.Color("#3f3f3f"))
)

type todayTaskDelegate struct {} 

func (d todayTaskDelegate) Height() int                             { return 3 }
func (d todayTaskDelegate) Spacing() int                            { return 1 }
func (d todayTaskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d todayTaskDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	t, ok := listItem.(task.Task)
	if !ok {
		return
	}

    priorityString := fmt.Sprintf("%d", t.Priority())

    renderedPriotity := unselectedPriorityStyle.Render(priorityString)
    renderedTitle := unselectedTitleStyle.Render(t.Title())
    renderedDescription := unselectedDescriptionStyle.Render(t.Description())

    if index == m.Index() {
        renderedPriotity = selectedPriorityStyle.Render(priorityString)
        renderedTitle = selectedTitleStyle.Width(m.Width()-8).Render(t.Title())
        renderedDescription = selectedDescriptionStyle.Width(m.Width()-8).Render(t.Description())
    } 

    item := lipgloss.JoinVertical(lipgloss.Left, renderedTitle, renderedDescription)
    item = lipgloss.JoinHorizontal(lipgloss.Top, renderedPriotity, item)

	fmt.Fprint(w, item)
}

