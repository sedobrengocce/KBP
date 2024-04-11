package board

import "strconv"

type errorColNumber struct{
    number int
}

func (e *errorColNumber) Error() string {
    return "Invalid column number: " + strconv.Itoa(e.number)
}

type emptyColumn struct{}

func (e *emptyColumn) Error() string {
    return "Column is empty"
}

type dbError struct {
    err error
}

func (e *dbError) Error() string {
    return e.err.Error()
}

