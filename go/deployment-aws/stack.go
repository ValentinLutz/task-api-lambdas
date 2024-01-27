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

func NewGetTasksFunction(stack awscdk.Stack, config *StageConfig) awslambda.Function {
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
		stack, jsii.String("V1GetTasks"), &awslambda.FunctionProps{
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-v1-get-tasks"),
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

func NewPostTasksFunction(stack awscdk.Stack, config *StageConfig) awslambda.Function {
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
		stack, jsii.String("V1PostTasks"), &awslambda.FunctionProps{
			Code: awslambda.Code_FromAsset(
				jsii.String("../lambda-v1-post-tasks"),
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
	getTasksFunction := NewGetTasksFunction(stack, config)
	postTasksFunction := NewPostTasksFunction(stack, config)

	openApiSpecs, err := template.ParseFiles("../api-definition/task-api-v1.yaml")
	if err != nil {
		panic(err)
	}

	var orderApiV1 bytes.Buffer
	err = openApiSpecs.Execute(
		&orderApiV1, map[string]string{
			"GetTasksFunctionArn":  *getTasksFunction.FunctionArn(),
			"PostTasksFunctionArn": *postTasksFunction.FunctionArn(),
		},
	)
	if err != nil {
		panic(err)
	}

	var orderApiV1Spec map[string]interface{}
	err = yaml.Unmarshal(orderApiV1.Bytes(), &orderApiV1Spec)
	if err != nil {
		panic(err)
	}

	restApi := awsapigateway.NewSpecRestApi(
		stack, NewIdWithStage(config, "TaskApi"), &awsapigateway.SpecRestApiProps{
			EndpointTypes: &[]awsapigateway.EndpointType{
				awsapigateway.EndpointType_REGIONAL,
			},
			ApiDefinition: awsapigateway.ApiDefinition_FromInline(orderApiV1Spec),
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
