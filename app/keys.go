package app

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Today, k.BoardList, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
        {k.Today, k.BoardList}, // first column
		{k.Help, k.Quit},                // second column
	}
}

type keyMap struct {
	Today   key.Binding
	BoardList   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Today: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("ctrl+t", "Today Board"),
	),
	BoardList: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctlr+l", "Board List"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

