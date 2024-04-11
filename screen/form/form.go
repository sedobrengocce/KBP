package form

import (
	"kaban-board-plus/common/component"
	"kaban-board-plus/component/task"
	msgs "kaban-board-plus/common/msg"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    checkBoxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9800"))
    checkBoxCheckedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9800")).Bold(true)
)

type formType int

const (
    newTask formType = iota
    newBoard
)

type Form struct {
    formType formType
    title textinput.Model
    description textarea.Model
    help help.Model
    priority task.Priority
    highlightedPriority task.Priority
    referrer component.Screen
}

func titleInput() textinput.Model {
    title := textinput.New()
    title.Placeholder = "Title"
    title.Prompt = "Title: "
    title.Focus()
    return title
}

func descriptionInput() textarea.Model {
    description := textarea.New()
    description.Placeholder = "Description"
    description.Prompt = "Description: "
    return description
}

func newHelp() help.Model {
    help := help.New()
    help.ShowAll = true
    return help
}

func radioButton(label string, highlighted bool, checked bool) string {
    if checked {
        return checkBoxCheckedStyle.Render("☑ " + label)
    } else if highlighted {
        return checkBoxStyle.Render("☐ " + label)
    } 
    return "☐ " + label
}

func NewNewTaskForm(s component.Screen) *Form {
    return &Form{
        formType: newTask,
        title: titleInput(),
        description: descriptionInput(),
        help: newHelp(),
        referrer: s,
    }
}

func NewNewBoardForm(s component.Screen) *Form {
    return &Form{
        formType: newBoard,
        title: titleInput(),
        help: newHelp(),
        referrer: s,
    }
}

func (f *Form) ToggleHelp() {
    f.help.ShowAll = !f.help.ShowAll
}

func (f Form) radioButtonsView() string {
    p := f.priority
    hp := f.highlightedPriority
    return lipgloss.JoinVertical(
        lipgloss.Center,
        radioButton("Low", hp == task.Low, p == task.Low),
        radioButton("Medium", hp == task.Medium, p == task.Medium),
        radioButton("High", hp == task.High, p == task.High),
    )
}

func (f Form) Init() tea.Cmd {
    return nil
}

func (f Form) updateRadioButtons(msg tea.Msg) (component.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter):
            return f.referrer, msgs.NewNewTaskMsg(f.title.Value(), f.description.Value(), int(f.priority))
        case key.Matches(msg, keys.Space):
            f.priority = f.highlightedPriority
        case key.Matches(msg, keys.Up):
            f.highlightedPriority = task.Priority((int(f.highlightedPriority) - 1) % 3)
        case key.Matches(msg, keys.Down):
            f.highlightedPriority = task.Priority((int(f.highlightedPriority) + 1) % 3)
        case key.Matches(msg, keys.Left):
            f.description.Focus()
            return &f, textarea.Blink
        }
    }
    return &f, nil
}

func (f Form) updateDescription(msg tea.Msg) (component.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter, keys.Right):
            f.description.Blur()
            return &f, nil
        case key.Matches(msg, keys.Left):
            f.description.Blur()
            f.title.Focus()
            return &f, textinput.Blink
        }
    }
    var cmd tea.Cmd
    f.description, cmd = f.description.Update(msg)
    return &f, cmd
}

func (f Form) updateTitle(msg tea.Msg) (component.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Enter):
            f.title.Blur()
            if f.formType == newTask {
                f.description.Focus()
                return &f, textarea.Blink
            }
            return f.referrer, msgs.NewNewBoardMsg(f.title.Value())
        case key.Matches(msg, keys.Right):
            if f.formType != newTask {
                return &f, nil
            }
            f.title.Blur()
            f.description.Focus()
            return &f, textarea.Blink
        }
    }
    var cmd tea.Cmd
    f.title, cmd = f.title.Update(msg)
    return &f, cmd
}

func (f Form) Update(msg tea.Msg) (component.Screen, tea.Cmd) {
    if keyMsg, ok := msg.(tea.KeyMsg); ok {
        if key.Matches(keyMsg, keys.Back) {
            return f.referrer.Update(nil)
        }
    }
    switch f.formType {
    case newTask:
        if(f.title.Focused()) {
            return f.updateTitle(msg)
        } else if(f.description.Focused()) {
            return f.updateDescription(msg)
        } else {
            return f.updateRadioButtons(msg)
        }
    case newBoard:
        return f.updateTitle(msg)
    }
    return &f, nil
}

func (f Form) View() string {
    if(f.formType == newTask) {
        return lipgloss.JoinVertical(
            lipgloss.Center, 
            f.title.View(), 
            f.description.View(), 
            f.radioButtonsView(),
            lipgloss.PlaceVertical(2, lipgloss.Bottom, f.help.View(keys)),
        )
    }
    return lipgloss.JoinVertical(
        lipgloss.Center,
        f.title.View(),
        lipgloss.PlaceVertical(2, lipgloss.Bottom, f.help.View(keys)),
    )
}

