package board

import (
	"fmt"
	"kaban-board-plus/component/task"
)

type colNumberError struct{
    number int
}

func NewColNumberError(n int) *colNumberError {
    return &colNumberError{number: n}
}

func (e *colNumberError) Error() string {
    return "Invalid column number: " + fmt.Sprint("%d", e.number)
}

type emptyColumn struct{}

func NewEmptyColumnError() *emptyColumn {
    return &emptyColumn{}
}

func (e *emptyColumn) Error() string {
    return "Column is empty"
}

type dbError struct {
    err error
}

func NewDbError(err error) *dbError {
    return &dbError{err: err}
}

func (e *dbError) Error() string {
    return e.err.Error()
}

type statusTypeError struct {
    status task.Status
}

func newStatusTypeError(s task.Status) *statusTypeError {
    return &statusTypeError{status: s}
}

func (e *statusTypeError) Error() string {
    return "Invalid status type: " + fmt.Sprint("%d", e.status)
}
