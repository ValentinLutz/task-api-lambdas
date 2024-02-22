package main

import (
	"bytes"
	"root/library-golang/cdk"
	"text/template"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)

func NewRestApi(stack awscdk.Stack, config *cdk.StageConfig, database awsrds.DatabaseInstance) awsapigateway.SpecRestApi {
	_ = cdk.NewGoFunction(stack, config, "DatabaseMigration", "../../service-migration/lambda", database)
	deleteTaskFunction := cdk.NewGoFunction(stack, config, "V1DeleteTask", "../lambda-v1-delete-task", database)
	getTaskFunction := cdk.NewGoFunction(stack, config, "V1GetTask", "../lambda-v1-get-task", database)
	getTasksFunction := cdk.NewGoFunction(stack, config, "V1GetTasks", "../lambda-v1-get-tasks", database)
	postTasksFunction := cdk.NewGoFunction(stack, config, "V1PostTasks", "../lambda-v1-post-tasks", database)
	putTaskFunction := cdk.NewGoFunction(stack, config, "V1PutTask", "../lambda-v1-put-task", database)

	openApiSpecs, err := template.ParseFiles("../../api-definition/task-api-v1.yaml")
	if err != nil {
		panic(err)
	}

	var orderApiV1 bytes.Buffer
	err = openApiSpecs.Execute(
		&orderApiV1, map[string]string{
			"DeleteTaskFunctionArn": *deleteTaskFunction.FunctionArn(),
			"GetTaskFunctionArn":    *getTaskFunction.FunctionArn(),
			"GetTasksFunctionArn":   *getTasksFunction.FunctionArn(),
			"PostTasksFunctionArn":  *postTasksFunction.FunctionArn(),
			"PutTaskFunctionArn":    *putTaskFunction.FunctionArn(),
		},
	)
	if err != nil {
		panic(err)
	}

	var apiV1Spec map[string]interface{}
	err = yaml.Unmarshal(orderApiV1.Bytes(), &apiV1Spec)
	if err != nil {
		panic(err)
	}

	restApi := awsapigateway.NewSpecRestApi(
		stack, cdk.NewIdWithStage(config, "TaskApi"), &awsapigateway.SpecRestApiProps{
			EndpointTypes: &[]awsapigateway.EndpointType{
				awsapigateway.EndpointType_REGIONAL,
			},
			ApiDefinition: awsapigateway.ApiDefinition_FromInline(apiV1Spec),
			DeployOptions: &awsapigateway.StageOptions{
				StageName: jsii.String(config.Environment),
			},
		},
	)

	deleteTaskFunction.AddPermission(
		jsii.String("AllowApiGatewayInvoke"), &awslambda.Permission{
			Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			SourceArn: restApi.ArnForExecuteApi(jsii.String("DELETE"), jsii.String("/v1/tasks/{task_id}"), nil),
		},
	)
	getTaskFunction.AddPermission(
		jsii.String("AllowApiGatewayInvoke"), &awslambda.Permission{
			Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			SourceArn: restApi.ArnForExecuteApi(jsii.String("GET"), jsii.String("/v1/tasks/{task_id}"), nil),
		},
	)
	getTasksFunction.AddPermission(
		jsii.String("AllowApiGatewayInvoke"), &awslambda.Permission{
			Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			SourceArn: restApi.ArnForExecuteApi(jsii.String("GET"), jsii.String("/v1/tasks"), nil),
		},
	)
	postTasksFunction.AddPermission(
		jsii.String("AllowApiGatewayInvoke"), &awslambda.Permission{
			Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			SourceArn: restApi.ArnForExecuteApi(jsii.String("POST"), jsii.String("/v1/tasks"), nil),
		},
	)
	putTaskFunction.AddPermission(
		jsii.String("AllowApiGatewayInvoke"), &awslambda.Permission{
			Principal: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
			SourceArn: restApi.ArnForExecuteApi(jsii.String("PUT"), jsii.String("/v1/tasks/{task_id}"), nil),
		},
	)

	return restApi
}

func NewStack(scope constructs.Construct, id *string, config *cdk.StageConfig) awscdk.Stack {
	stack := awscdk.NewStack(
		scope, id, &awscdk.StackProps{Env: &awscdk.Environment{Account: &config.Account, Region: &config.Region}},
	)

	vpc := cdk.NewVpc(stack)
	database := cdk.NewDatabase(stack, vpc)
	cdk.NewVpcEndpoint(stack, vpc)
	NewRestApi(stack, config, database)

	return stack
}
