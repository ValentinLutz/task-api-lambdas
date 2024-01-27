package main

import (
	"bytes"
	"text/template"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"gopkg.in/yaml.v3"
)

func NewFunction(stack awscdk.Stack, config *StageConfig, functionName string, bootstrapPath string) awslambda.Function {
	env := map[string]*string{
		"DB_HOST":      &config.databaseProps.host,
		"DB_PORT":      &config.databaseProps.port,
		"DB_NAME":      &config.databaseProps.name,
		"DB_SECRET_ID": &config.databaseProps.secret,
	}
	if config.endpointUrl != nil {
		env["AWS_ENDPOINT_URL"] = config.endpointUrl
	}

	lambdaFunction := awslambda.NewFunction(
		stack, jsii.String(functionName), &awslambda.FunctionProps{
			Code: awslambda.Code_FromAsset(
				jsii.String(bootstrapPath),
				&awss3assets.AssetOptions{
					IgnoreMode: awscdk.IgnoreMode_GIT,
					Exclude: &[]*string{
						jsii.String("**"),
						jsii.String("!bootstrap"),
					},
				},
			),
			Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
			MemorySize:   jsii.Number(128),
			Handler:      jsii.String("bootstrap"),
			Architecture: config.lambdaConfig.architecture,
			Environment:  &env,
		},
	)

	return lambdaFunction
}

func NewRestApi(stack awscdk.Stack, config *StageConfig) awscdk.Stack {
	getTaskFunction := NewFunction(stack, config, "V1GetTask", "../lambda-v1-get-task")
	getTasksFunction := NewFunction(stack, config, "V1GetTasks", "../lambda-v1-get-tasks")
	postTasksFunction := NewFunction(stack, config, "V1PostTasks", "../lambda-v1-post-tasks")

	openApiSpecs, err := template.ParseFiles("../api-definition/task-api-v1.yaml")
	if err != nil {
		panic(err)
	}

	var orderApiV1 bytes.Buffer
	err = openApiSpecs.Execute(
		&orderApiV1, map[string]string{
			"GetTaskFunctionArn":   *getTaskFunction.FunctionArn(),
			"GetTasksFunctionArn":  *getTasksFunction.FunctionArn(),
			"PostTasksFunctionArn": *postTasksFunction.FunctionArn(),
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
		stack, NewIdWithStage(config, "TaskApi"), &awsapigateway.SpecRestApiProps{
			EndpointTypes: &[]awsapigateway.EndpointType{
				awsapigateway.EndpointType_REGIONAL,
			},
			ApiDefinition: awsapigateway.ApiDefinition_FromInline(apiV1Spec),
			DeployOptions: &awsapigateway.StageOptions{
				StageName: jsii.String(config.environment),
			},
			Policy: awsiam.NewPolicyDocument(
				&awsiam.PolicyDocumentProps{
					Statements: &[]awsiam.PolicyStatement{
						awsiam.NewPolicyStatement(
							&awsiam.PolicyStatementProps{
								Effect: awsiam.Effect_ALLOW,
								Actions: &[]*string{
									jsii.String("execute-api:Invoke"),
								},
								Resources: &[]*string{
									jsii.String("*"),
								},
								Principals: &[]awsiam.IPrincipal{
									awsiam.NewAnyPrincipal(),
								},
							},
						),
					},
				},
			),
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

	return stack
}

func NewStack(scope constructs.Construct, id *string, config *StageConfig) awscdk.Stack {
	stack := awscdk.NewStack(
		scope, id, &awscdk.StackProps{Env: &awscdk.Environment{Account: &config.account, Region: &config.region}},
	)

	return NewRestApi(stack, config)
}
