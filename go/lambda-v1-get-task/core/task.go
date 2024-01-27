package core

import (
	"context"
	"fmt"
	"root/go/lambda-v1-get-task/outgoing"

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

func (service *TaskService) GetTask(ctx context.Context, taskId uuid.UUID) (outgoing.TaskEntity, error) {
	taskEntity, err := service.taskRepository.FindByTaskId(ctx, taskId)
	if err != nil {
		return outgoing.TaskEntity{}, fmt.Errorf("failed to get task by task id: %w", err)
	}

	return taskEntity, nil
}
