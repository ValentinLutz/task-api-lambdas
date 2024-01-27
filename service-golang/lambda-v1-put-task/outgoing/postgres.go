package outgoing

import (
	"context"
	"fmt"

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

func (taskRepository *TaskRepository) Update(ctx context.Context, taskEntity TaskEntity) error {
	result, err := taskRepository.database.NamedExecContext(
		ctx,
		`
		UPDATE public.tasks
		SET title = :title, description = :description
		WHERE task_id = :task_id`,
		taskEntity,
	)
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
