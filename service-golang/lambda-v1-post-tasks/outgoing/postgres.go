package outgoing

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	database *sqlx.DB
}

func NewTaskRepository(database *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		database: database,
	}
}

func (taskRepository *TaskRepository) Save(ctx context.Context, taskEntity TaskEntity) error {
	_, err := taskRepository.database.NamedExecContext(
		ctx,
		"INSERT INTO public.tasks (task_id, title, description) VALUES (:task_id, :title, :description)",
		taskEntity,
	)
	if err != nil {
		return err
	}

	return nil
}
