package core

import (
	"context"
	"fmt"
	"root/go/lambda-v1-get-tasks/outgoing"
)

type TaskService struct {
	taskRepository *outgoing.TaskRepository
}

func NewTaskService(orderRepository *outgoing.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: orderRepository,
	}
}

func (service *TaskService) GetTasks(ctx context.Context) ([]outgoing.TaskEntity, error) {
	tasks, err := service.taskRepository.FindAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}
