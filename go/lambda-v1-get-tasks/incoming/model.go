package incoming

import (
	"root/go/lambda-v1-get-tasks/outgoing"
)

//go:generate oapi-codegen --config ../../api-definition/oapi-codgen.yaml ../../api-definition/task-api-v1.yaml

func NewTasksResponse(tasks []outgoing.TaskEntity) TasksResponse {
	tasksResponse := make([]TaskResponse, 0)
	for _, task := range tasks {
		tasksResponse = append(tasksResponse, NewTaskResponse(task))
	}

	return tasksResponse
}

func NewTaskResponse(task outgoing.TaskEntity) TaskResponse {
	return TaskResponse{
		Description: task.Description,
		Title:       task.Title,
		TaskId:      task.TaskId,
	}
}
