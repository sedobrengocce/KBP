package task

import componet "kaban-board-plus/common/component"

type Status int

type Priority int

const (
    Todo Status = iota
    InProgress
    Done
)

const (
    High Priority = iota + 1
    Medium
    Low
)

type Task struct {
    id int
    name string
    description string
    priority Priority
    isArchived bool
    completedDate componet.CompletedDate
    status Status
}

func NewTask(id int, title, description string, priority Priority, status Status, isArchived bool) Task {
    task :=Task{
        id: id,
        name: title,
        description: description,
        priority: priority,
        status: status,
        isArchived: isArchived,
    }

    return task
}

func (t Task) FilterValue() string {
    return t.name + t.description
}

func (t Task) Title() string {
    return t.name
}

func (t Task) ID() int {
    return t.id
}

func (t Task) Description() string {
    return t.description
}

func (t Task) Priority() Priority {
    return t.priority
}

func (t Task) Status() Status {
    return t.status
}

func (t Task) IsArchived() bool {
    return t.isArchived
}

func (t *Task) SetPriority(priority Priority) {
    t.priority = priority
}

func (t *Task) ToggleArchive() {
    t.isArchived = !t.isArchived
}

func (t *Task) SetDescription(description string) {
    t.description = description
}

func (t *Task) SetTitle(name string) {
    t.name = name
}

