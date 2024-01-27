package outgoing

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	database *sqlx.DB
}

func NewTaskRepository(database *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		database: database,
	}
}

func (taskRepository *TaskRepository) FindByTaskId(ctx context.Context, taskId uuid.UUID) (TaskEntity, error) {
	var taskEntity TaskEntity
	err := taskRepository.database.GetContext(
		ctx,
		&taskEntity,
		"SELECT task_id, title, description FROM public.tasks WHERE task_id = $1",
		taskId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return TaskEntity{}, ErrTaskNotFound
	}
	if err != nil {
		return TaskEntity{}, err
	}

	return taskEntity, nil
}
