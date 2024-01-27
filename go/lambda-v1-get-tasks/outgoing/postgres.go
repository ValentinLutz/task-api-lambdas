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

func (taskRepository *TaskRepository) FindAllTasks(ctx context.Context) ([]TaskEntity, error) {
	var taskEntities []TaskEntity
	err := taskRepository.database.SelectContext(ctx, &taskEntities, "SELECT task_id, title, description FROM public.tasks")
	if err != nil {
		return nil, err
	}
	return taskEntities, nil
}
