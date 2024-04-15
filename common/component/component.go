package component

import (
	"log"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
    Init() tea.Cmd
    Update(msg tea.Msg) (Screen, tea.Cmd)
    View() string
    ToggleHelp()
    SetSize(width, height int)
}

type CompletedDate struct {
    day int
    month int
    year int
}

func newCompletedDate(day, month, year int) *CompletedDate {
    return &CompletedDate{
        day: day,
        month: month,
        year: year,
    }
}

func CompletedDateFromDateString(dateString string) (*CompletedDate, error) {
    date, err := time.Parse(time.DateOnly, dateString)
    if err != nil {
        log.Print("Invalid date string: ", err)
        return nil, err
    }

    year, month, day := date.Date()

    cd := newCompletedDate(day, int(month), year)

    return cd, nil
}

func CompletedDateFromTime(t time.Time) *CompletedDate {
    year, month, day := t.Date()

    return newCompletedDate(day, int(month), year)
}

func (cd *CompletedDate) GetDateString() string {
    dateString := strconv.Itoa(cd.year) + "-" + strconv.Itoa(cd.month) + "-" + strconv.Itoa(cd.day)
    return dateString
}

func (cd *CompletedDate) GetDate() (time.Time, error) {
    dateString := cd.GetDateString()
    return time.Parse(time.DateOnly, dateString)
}
