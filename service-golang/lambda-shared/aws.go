package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"engine"`
}

func GetDatabaseSecret(config aws.Config) (DatabaseSecret, error) {
	secretId, ok := os.LookupEnv("DB_SECRET_ID")
	if !ok {
		return DatabaseSecret{}, fmt.Errorf("env DB_SECRET_ID not set")
	}

	client := secretsmanager.NewFromConfig(config)
	secretValue, err := client.GetSecretValue(
		context.Background(), &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretId),
		},
	)
	if err != nil {
		return DatabaseSecret{}, fmt.Errorf("failed to get secret with id %s: %w", secretId, err)
	}

	var secret DatabaseSecret
	err = json.Unmarshal([]byte(*secretValue.SecretString), &secret)
	if err != nil {
		return DatabaseSecret{}, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return secret, nil
}
