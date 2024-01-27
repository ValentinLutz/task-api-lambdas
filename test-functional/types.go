package testfunctional

import "github.com/google/uuid"

type TaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

type TaskResponse struct {
	TaskId      uuid.UUID `json:"task_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
}

type TasksResponse = []TaskResponse
