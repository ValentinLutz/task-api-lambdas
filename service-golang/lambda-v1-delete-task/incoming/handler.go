package incoming

import (
	"context"
	"errors"
	"fmt"
	"os"
	"root/service-golang/lambda-shared"
	"root/service-golang/lambda-v1-delete-task/core"
	"root/service-golang/lambda-v1-delete-task/outgoing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Database    *sqlx.DB
	TaskService *core.TaskService
}

func NewHandler() (*Handler, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws default config: %w", err)
	}

	secret, err := shared.GetSecret(cfg, os.Getenv("DB_SECRET_ID"))
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
		Database:    database,
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

	err = handler.TaskService.DeleteTask(ctx, taskId)
	if errors.Is(err, outgoing.ErrTaskNotFound) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	} else if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to delete task: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}