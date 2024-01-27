package library

import (
	"fmt"
	"os"
	"runtime"
)

type DatabaseProps struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func GetOrSetDefaultDatabaseEnvVars() *DatabaseProps {
	return &DatabaseProps{
		Host:     GetValueOrSetDefault("DB_HOST", "localhost"),
		Port:     GetValueOrSetDefault("DB_PORT", "5432"),
		Name:     GetValueOrSetDefault("DB_NAME", "test"),
		User:     GetValueOrSetDefault("DB_USER", "test"),
		Password: GetValueOrSetDefault("DB_PASS", "test"),
	}
}

type BuildProps struct {
	OperatingSystem string
	Architecture    string
}

func GetOrSetDefaultBuildEnvVars() *BuildProps {
	stageEnvVars := GetOrSetDefaultStageEnvVars()

	GetValueOrSetDefault("CGO_ENABLED", "0")

	if stageEnvVars.Environment == "local" {
		return &BuildProps{
			OperatingSystem: GetValueOrSetDefault("GOOS", runtime.GOOS),
			Architecture:    GetValueOrSetDefault("GOARCH", runtime.GOARCH),
		}
	}
	return &BuildProps{
		OperatingSystem: GetValueOrSetDefault("GOOS", "linux"),
		Architecture:    GetValueOrSetDefault("GOARCH", "arm64"),
	}
}

type StageProps struct {
	Environment string
	Region      string
}

func GetOrSetDefaultStageEnvVars() *StageProps {
	return &StageProps{
		Environment: GetValueOrSetDefault("ENVIRONMENT", "local"),
		Region:      GetValueOrSetDefault("REGION", "eu-central-1"),
	}
}

func GetValueOrSetDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("env '%s' not set, defaulting to '%s'\n", key, defaultValue)
		err := os.Setenv(key, defaultValue)
		if err != nil {
			panic(err)
		}
		return defaultValue
	}
	return value
}
