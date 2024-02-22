package main

import (
	"root/library-golang/cdk"

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
)

var stageConfigs = map[string]*cdk.StageConfig{
	"eu-central-1-test": {
		Account:     "489721517942",
		Region:      "eu-central-1",
		Environment: "test",
		LambdaConfig: cdk.LambdaConfig{
			Architecture: awslambda.Architecture_ARM_64(),
		},
	},
}
