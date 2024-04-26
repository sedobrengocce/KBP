package board

import "github.com/charmbracelet/bubbles/key"

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Action, k.NewTask}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Action, k.NewTask},                // second column
	}
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Left   key.Binding
	Action  key.Binding
    NewTask key.Binding
    Archive key.Binding
    ArchiveAll key.Binding
    ShowHideArchive key.Binding
    Delete  key.Binding
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
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Action: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "action"),
	),
    NewTask: key.NewBinding(
        key.WithKeys("n"),
        key.WithHelp("n", "new task (not in today board)"),
    ),
	Archive: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "toggle archive"),
	),
    ShowHideArchive: key.NewBinding(
        key.WithKeys("ctrl+a"),
    ),
    ArchiveAll: key.NewBinding(
        key.WithKeys("shift+a"),
    ),
    Delete: key.NewBinding(
        key.WithKeys("d"),
    ),
}

