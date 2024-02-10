package incoming

import (
	"context"
	"encoding/json"
	"fmt"
	shared "root/service-golang/lambda-shared"
	"root/service-golang/lambda-v1-get-tasks/core"
	"root/service-golang/lambda-v1-get-tasks/outgoing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
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

func (handler *Handler) Invoke(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	taskEntities, err := handler.TaskService.GetTasks(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to get taskEntities: %w", err)
	}

	tasksResponse := NewTasksResponse(taskEntities)
	tasksResponseBody, err := json.Marshal(tasksResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to marshal taskEntities: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(tasksResponseBody),
	}, nil
}
