// TODO: Better style for Today task and fro Archived task
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
    unselectedPPriorityStyle = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).Width(3).Height(3).
                        Foreground(lipgloss.Color("#ff8550")).
                        BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff8550"))
    unselectedPTitleStyle = lipgloss.NewStyle().Height(2).
                        Foreground(lipgloss.Color("#ff8550"))
    unselectedPDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8f50"))
)

var (
    selectedPPriorityStyle = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).Width(3).Height(3).
                        Bold(true).Foreground(lipgloss.Color("#ff9800")).Background(lipgloss.Color("#3f3f3f")).
                        BorderStyle(lipgloss.ThickBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff9800")).BorderBackground(lipgloss.Color("#3f3f3f"))
    selectedPTitleStyle = lipgloss.NewStyle().Height(2).
                        Bold(true).Foreground(lipgloss.Color("#ff9800")).Background(lipgloss.Color("#3f3f3f"))
    selectedPDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff8f50")).Background(lipgloss.Color("#3f3f3f"))
)

var (
    selectedToday = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).PaddingLeft(2).Width(5).Height(3).
                        Bold(true).Foreground(lipgloss.Color("#ff9800")).Background(lipgloss.Color("#3f3f3f")).
                        BorderStyle(lipgloss.ThickBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff9800")).BorderBackground(lipgloss.Color("#3f3f3f"))
    unselectedToday = lipgloss.NewStyle().Align(lipgloss.Center,lipgloss.Center).PaddingLeft(2).Width(5).Height(3).
                        Foreground(lipgloss.Color("#ff8550")).
                        BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(lipgloss.Color("#ff8550"))
)

type projectTaskDelegate struct {} 

func (d projectTaskDelegate) Height() int                             { return 3 }
func (d projectTaskDelegate) Spacing() int                            { return 1 }
func (d projectTaskDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d projectTaskDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
    t, ok := listItem.(task.Task)
    if !ok {
        return
    }

    priorityString := fmt.Sprintf("%d", t.Priority())
    todayString := "T"
    archivedString := "A"

    renderedPriority := unselectedPPriorityStyle.Render(priorityString)
    renderedTitle := unselectedPTitleStyle.Render(t.Title())
    renderedDescription := unselectedPDescriptionStyle.Render(t.Description())
    renderedTodady := unselectedPPriorityStyle.Render(todayString)
    renderedArchived := unselectedPPriorityStyle.Render(archivedString)

    if index == m.Index() {
        margin := 8;
        if t.IsToday() {
            margin += 4
        }
        if t.IsArchived() {
            margin += 4
        }
        renderedPriority = selectedPPriorityStyle.Render(priorityString)
        renderedTitle = selectedPTitleStyle.Width(m.Width()-margin).Render(t.Title())
        renderedDescription = selectedPDescriptionStyle.Width(m.Width()-margin).Render(t.Description())
        renderedTodady = selectedPPriorityStyle.Render(todayString)
        renderedArchived = selectedPPriorityStyle.Render(archivedString)
    } 

    item := lipgloss.JoinVertical(lipgloss.Left, renderedTitle, renderedDescription)
    item = lipgloss.JoinHorizontal(lipgloss.Top, renderedPriority, item)
    if t.IsToday() {
    item = lipgloss.JoinHorizontal(lipgloss.Top, renderedTodady, item)
    } 
    if t.IsArchived() {
    item = lipgloss.JoinHorizontal(lipgloss.Top, renderedArchived, item)
    }
    fmt.Fprint(w, item)
}

