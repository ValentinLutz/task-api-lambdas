package outgoing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrTaskNotFound = fmt.Errorf("task not found")

type TaskRepository struct {
	database *sqlx.DB
}

func NewTaskRepository(database *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		database: database,
	}
}

func (taskRepository *TaskRepository) DeleteByTaskId(ctx context.Context, taskId uuid.UUID) error {
	result, err := taskRepository.database.ExecContext(ctx, "DELETE FROM public.tasks WHERE task_id = $1", taskId)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return ErrTaskNotFound
	}

	return nil

}
