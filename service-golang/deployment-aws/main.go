package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	config, err := NewStageConfig()
	if err != nil {
		panic(err)
	}

	app := awscdk.NewApp(nil)
	tags := awscdk.Tags_Of(app)
	tags.Add(jsii.String("custom:region"), &config.region, &awscdk.TagProps{})
	tags.Add(jsii.String("custom:environment"), &config.environment, &awscdk.TagProps{})

	NewStack(app, NewIdWithStage(config, "TaskResource"), config)

	app.Synth(nil)
}
