package list

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.NewBoard, k.DeleteBoard}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down}, // first column
		{k.Enter, k.NewBoard, k.DeleteBoard},                // second column
	}
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
    NewBoard key.Binding
    DeleteBoard key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select board"),
	),
    NewBoard: key.NewBinding(
        key.WithKeys("n"),
        key.WithHelp("n", "new board"),
    ),
    DeleteBoard: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "delete board"),
    ),
}

