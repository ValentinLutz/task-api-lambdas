package core

import (
	"context"
	"fmt"
	"root/service-golang/lambda-v1-get-tasks/outgoing"
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
	tasksEntities, err := service.taskRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasksEntities: %w", err)
	}

	return tasksEntities, nil
}
