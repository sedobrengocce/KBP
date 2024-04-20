package list

type boardElement struct {
    id int
    name string
}

func (b boardElement) FilterValue() string {
    return b.name
}

func newBoardElement(id int, name string) *boardElement {
    return &boardElement{
        id: id,
        name: name,
    }
}

func (b boardElement) title() string {
    return b.name
}

