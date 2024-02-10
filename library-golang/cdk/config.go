package cdk

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type StageConfig struct {
	Account      string
	Region       string
	Environment  string
	EndpointUrl  *string
	LambdaConfig LambdaConfig
}

type LambdaConfig struct {
	Architecture awslambda.Architecture
}

var (
	ErrStageConfigNotFound      = fmt.Errorf("stage config not found")
	ErrRegionEnvNotSet          = fmt.Errorf("env REGION not set")
	ErrEnvironmentEnvNotSet     = fmt.Errorf("env ENVIRONMENT not set")
	ErrArchitectureNotSupported = fmt.Errorf("architecture not supported")
)

func NewStageConfig(stageConfigs map[string]*StageConfig) (*StageConfig, error) {
	region, ok := os.LookupEnv("REGION")
	if !ok {
		return nil, ErrRegionEnvNotSet
	}

	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return nil, ErrEnvironmentEnvNotSet
	}

	stageKey := region + "-" + env
	stage, ok := stageConfigs[stageKey]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrStageConfigNotFound, stageKey)
	}

	return stage, nil
}

func NewIdWithStage(stage *StageConfig, id string) *string {
	envTitleFormat := cases.Title(language.English).String(stage.Environment)
	return jsii.String(id + envTitleFormat)
}

func GetArchitecture() awslambda.Architecture {
	switch runtime.GOARCH {
	case "amd64":
		return awslambda.Architecture_X86_64()
	case "arm64":
		return awslambda.Architecture_ARM_64()
	default:
		panic(fmt.Errorf("%w: %s", ErrArchitectureNotSupported, runtime.GOARCH))
	}
}
