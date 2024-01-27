package core

import (
	"context"
	"fmt"
	"root/service-golang/lambda-v1-put-task/outgoing"

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

func (service *TaskService) UpdateTask(ctx context.Context, taskId uuid.UUID, title string, description *string) error {
	taskEntity := outgoing.TaskEntity{
		TaskId:      taskId,
		Title:       title,
		Description: description,
	}
	err := service.taskRepository.Update(ctx, taskEntity)
	if err != nil {
		return fmt.Errorf("failed to update task by task id: %w", err)
	}

	return nil
}
