package core

import (
	"context"
	"fmt"
	"root/go/lambda-v1-post-tasks/outgoing"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepository *outgoing.TaskRepository
}

func NewTaskService(orderRepository *outgoing.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: orderRepository,
	}
}

func (service *TaskService) SaveTask(ctx context.Context, title string, description *string) (outgoing.TaskEntity, error) {
	taskId, err := uuid.NewV7()
	if err != nil {
		return outgoing.TaskEntity{}, fmt.Errorf("failed to generate task id: %w", err)
	}

	taskEntity := outgoing.TaskEntity{
		TaskId:      taskId,
		Title:       title,
		Description: description,
	}

	err = service.taskRepository.Save(ctx, taskEntity)
	if err != nil {
		return outgoing.TaskEntity{}, fmt.Errorf("failed to save task: %w", err)
	}

	return taskEntity, nil
}
