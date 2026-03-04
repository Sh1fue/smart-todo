package task

type Priotity int

type Status int
type Task struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priotity `json:"priority"`
}
