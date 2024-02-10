package incoming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	shared "root/service-golang/lambda-shared"
	"root/service-golang/lambda-v1-get-task/core"
	"root/service-golang/lambda-v1-get-task/outgoing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
)

type Handler struct {
	TaskService *core.TaskService
}

func NewHandler() (*Handler, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws default config: %w", err)
	}

	secret, err := shared.GetDatabaseSecret(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get database secret: %w", err)
	}

	dbConfig, err := shared.NewDatabaseConfig(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create database config: %w", err)
	}

	database, err := shared.NewDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	taskRepository := outgoing.NewTaskRepository(database)
	taskService := core.NewTaskService(taskRepository)

	return &Handler{
		TaskService: taskService,
	}, nil
}

func (handler *Handler) Invoke(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	taskIdString, ok := r.PathParameters["task_id"]
	if !ok {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("task id not found in path parameters")
	}

	taskId, err := uuid.Parse(taskIdString)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to parse task id: %w", err)
	}

	taskEntity, err := handler.TaskService.GetTask(ctx, taskId)
	if errors.Is(err, outgoing.ErrTaskNotFound) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to get task: %w", err)
	}

	tasksResponse := NewTaskResponse(taskEntity)
	tasksResponseBody, err := json.Marshal(tasksResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to marshal task: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(tasksResponseBody),
	}, nil
}
