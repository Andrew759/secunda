package enum

type Status int

// Статусы задач
const (
	Draft Status = iota + 1
	Todo
	InProgress
	Done
)

func (s Status) IsValid() bool {
	return s >= Draft && s <= Done
}
