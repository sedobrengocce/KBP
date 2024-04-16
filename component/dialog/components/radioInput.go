package dialogComponent

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
    focusedItemStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#DD7600")).Background(lipgloss.Color("#5C5C5C"))
    unfocusedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EFEFEF"))
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("Up", "k"),
	),
	Down: key.NewBinding(
        key.WithKeys("Down", "j"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
}

type RadioItem[T any] struct {
    title string
    value T
}

type radioInput[T any] struct {
    prompt string
    items []RadioItem[T]
    selectedItem int
    focusedItem int
    hasSelectedItem bool
    hasFocus bool
}

func NewRadioItem[T any](title string, value T) *RadioItem[T] {
    return &RadioItem[T]{
        title: title,
        value: value,
    }
}

func NewRadioInput[T any](prompt string, items []RadioItem[T]) *radioInput[T] {
    return &radioInput[T]{
        prompt: prompt,
        items: items,
        hasSelectedItem: false,
        hasFocus: false,
        focusedItem: 0,
        selectedItem: 0,
    }
}

func (ri *radioInput[T]) Focus() tea.Cmd {
    ri.hasFocus = true
    return nil
}

func (ri *radioInput[T]) Blur() tea.Cmd {
    ri.hasFocus = false 
    return nil
}

func (ri radioInput[T]) GetValue() T {
    return ri.items[ri.selectedItem].value
}

func (ri *radioInput[T]) nextItem() {
    ri.focusedItem = (ri.focusedItem + 1) % len(ri.items)
}

func (ri *radioInput[T]) prevItem() {
    numItems := len(ri.items)
    ri.focusedItem = (ri.focusedItem + 1 + numItems) % numItems
}

func (ri *radioInput[T]) selectItem() {
    if ri.selectedItem == ri.focusedItem {
        ri.hasSelectedItem = false
        return
    }
    ri.hasSelectedItem = true
    ri.selectedItem = ri.focusedItem
}

func (ri *radioInput[T]) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, keys.Up):
            ri.prevItem()
            return nil
        case key.Matches(msg, keys.Down):
            ri.nextItem()
            return nil
        case key.Matches(msg, keys.Enter):
            ri.selectItem()
            return nil
        }
    }   
    return nil
}

func (ri radioInput[T]) Render() string {
    checked := "[X] "
    unchecked := "[ ] "

    lines := []string{
        ri.prompt,
    }

    for i, e := range ri.items {
        var line string
        if i == ri.selectedItem && ri.hasSelectedItem {
            line = checked + e.title
        } else {
            line = unchecked + e.title
        }
        if i == ri.focusedItem && ri.hasFocus{
            lines = append(lines, focusedItemStyle.Render(line))
        } else {
            lines = append(lines, unfocusedItemStyle.Render(line))
        }
    }

    return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (ri radioInput[T]) IsFocused() bool {
    return ri.hasFocus
}

    
