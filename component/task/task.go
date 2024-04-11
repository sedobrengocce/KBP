package task

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
    dueDate string
    status Status
}

func NewTask(id int, title, description string, priority Priority) Task {
    task :=Task{
        id: id,
        name: title,
        description: description,
        priority: priority,
        status: Todo,
    }

    return task
}

func (t Task) WithDueDate(date string) Task {
    t.dueDate = date
    return t
}

func (t Task) WithStatus(status Status) Task {
    t.status = status
    return t
}

func (t Task) WithArchived(isArchived bool) Task {
    t.isArchived = isArchived
    return t
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

func (t Task) DueDate() string {
    return t.dueDate
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

