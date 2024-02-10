package shared

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

var (
	ErrSecretUsernameNotSet = fmt.Errorf("secret username not set")
	ErrSecretPasswordNotSet = fmt.Errorf("secret password not set")
	ErrDbHostEnvNotSet      = fmt.Errorf("secret host not set")
	ErrDbPortEnvNotSet      = fmt.Errorf("secret port not set")
	ErrDbNameEnvNotSet      = fmt.Errorf("secret name not set")
)

func NewDatabaseConfig(secret DatabaseSecret) (*DatabaseConfig, error) {
	if secret.Username == "" {
		return nil, ErrSecretUsernameNotSet
	}
	if secret.Password == "" {
		return nil, ErrSecretPasswordNotSet
	}

	if secret.Host == "" {
		return nil, ErrDbHostEnvNotSet
	}

	if secret.Port == 0 {
		return nil, ErrDbPortEnvNotSet
	}

	if secret.Name == "" {
		return nil, ErrDbNameEnvNotSet
	}

	return &DatabaseConfig{
		Host:     secret.Host,
		Port:     secret.Port,
		Name:     secret.Name,
		User:     secret.Username,
		Password: secret.Password,
	}, nil
}

func NewDatabase(config *DatabaseConfig) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s",
		config.Host, config.Port, config.Name, config.User, config.Password,
	)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
