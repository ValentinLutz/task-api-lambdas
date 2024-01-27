package outgoing

import "github.com/google/uuid"

type TaskEntity struct {
	TaskId      uuid.UUID `db:"task_id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
}
