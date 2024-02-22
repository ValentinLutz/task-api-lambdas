package incoming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	shared "root/service-golang/lambda-shared"
	"root/service-golang/lambda-v1-put-task/core"
	"root/service-golang/lambda-v1-put-task/outgoing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
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

	awsv2.AWSV2Instrumentor(&cfg.APIOptions)

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

	var taskRequest TaskRequest
	err = json.Unmarshal([]byte(r.Body), &taskRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to unmarshal task request: %w", err)
	}

	err = handler.TaskService.UpdateTask(ctx, taskId, taskRequest.Title, taskRequest.Description)
	if errors.Is(err, outgoing.ErrTaskNotFound) {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	}
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to update task: %w", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}
