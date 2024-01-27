package core

import (
	"context"
	"fmt"
	"root/go/lambda-v1-delete-task/outgoing"

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

func (service *TaskService) DeleteTask(ctx context.Context, taskId uuid.UUID) error {
	err := service.taskRepository.DeleteByTaskId(ctx, taskId)
	if err != nil {
		return fmt.Errorf("failed to delete task by task id: %w", err)
	}

	return nil
}
