package main

import (
	"context"
	"fmt"
	shared "root/service-golang/lambda-shared"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	Database *sqlx.DB
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

	return &Handler{
		Database: database,
	}, nil
}

func (handler *Handler) Invoke(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	query := `
	DROP TABLE IF EXISTS public.tasks;
	CREATE TABLE IF NOT EXISTS public.tasks
	(
		task_id         UUID    NOT NULL UNIQUE,
		title           TEXT,
		description     TEXT,
		PRIMARY KEY (task_id)
	);`

	_, err := handler.Database.Exec(query)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
