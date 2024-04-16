package dialogComponent

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Space key.Binding
    Enter key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("Up", "k"),
	),
	Down: key.NewBinding(
        key.WithKeys("Down", "j"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
	),
    Enter: key.NewBinding(
        key.WithKeys("enter"),
    ),
}

