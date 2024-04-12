package dialog

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Left, k.Right, k.Enter}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Right}, // first column
		{k.Enter},                // second column
	}
}

type keyMap struct {
	Left     key.Binding
	Right   key.Binding
	Enter  key.Binding
}

var keys = keyMap{
	Left: key.NewBinding(
		key.WithKeys("Left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
        key.WithKeys("Right", "l"),
        key.WithHelp("→/l", "move right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select board"),
	),
}


