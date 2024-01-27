package incoming

import (
	"root/go/lambda-v1-get-task/outgoing"
)

//go:generate oapi-codegen --config ../../api-definition/oapi-codgen.yaml ../../api-definition/task-api-v1.yaml

func NewTaskResponse(task outgoing.TaskEntity) TaskResponse {
	return TaskResponse{
		Description: task.Description,
		Title:       task.Title,
		TaskId:      task.TaskId,
	}
}
