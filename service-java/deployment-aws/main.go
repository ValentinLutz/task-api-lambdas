package main

import (
	"root/library-golang/cdk"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	config, err := cdk.NewStageConfig(stageConfigs)
	if err != nil {
		panic(err)
	}

	app := awscdk.NewApp(nil)
	tags := awscdk.Tags_Of(app)
	tags.Add(jsii.String("custom:region"), &config.Region, &awscdk.TagProps{})
	tags.Add(jsii.String("custom:environment"), &config.Environment, &awscdk.TagProps{})
	tags.Add(jsii.String("custom:language"), jsii.String("java"), &awscdk.TagProps{})

	NewStack(app, cdk.NewIdWithStage(config, "TaskResource"), config)

	app.Synth(nil)
}
