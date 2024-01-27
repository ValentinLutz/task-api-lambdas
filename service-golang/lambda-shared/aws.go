package shared

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSecret(config aws.Config, secretId string) (Secret, error) {
	conn := secretsmanager.NewFromConfig(config)

	secretValue, err := conn.GetSecretValue(
		context.Background(), &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretId),
		},
	)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to get secret with id %s: %w", secretId, err)
	}

	var secret Secret
	err = json.Unmarshal([]byte(*secretValue.SecretString), &secret)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return secret, nil
}
