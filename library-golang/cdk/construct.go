package cdk

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/jsii-runtime-go"
)

func NewJavaFunction(
	stack awscdk.Stack,
	config *StageConfig,
	functionName string,
	bootstrapPath string,
	database awsrds.DatabaseInstance,
) awslambda.Function {
	env := map[string]*string{
		"DB_SECRET_ID": database.Secret().SecretName(),
	}

	function := awslambda.NewFunction(
		stack, jsii.String(functionName), &awslambda.FunctionProps{
			Vpc:          database.Vpc(),
			Code:         awslambda.Code_FromAsset(jsii.String(bootstrapPath), &awss3assets.AssetOptions{}),
			Runtime:      awslambda.Runtime_JAVA_21(),
			MemorySize:   jsii.Number(128),
			Handler:      jsii.String("science.monke.incoming.Handler"),
			Architecture: config.LambdaConfig.Architecture,
			Environment:  &env,
		},
	)

	database.Secret().GrantRead(function, nil)
	database.Connections().AllowFrom(
		function,
		awsec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow access from java lambda to database"),
	)

	return function
}

func NewGoFunction(
	stack awscdk.Stack,
	config *StageConfig,
	functionName string,
	bootstrapPath string,
	database awsrds.DatabaseInstance,
) awslambda.Function {
	env := map[string]*string{
		"DB_SECRET_ID": database.Secret().SecretName(),
	}

	function := awslambda.NewFunction(
		stack, jsii.String(functionName), &awslambda.FunctionProps{
			Vpc: database.Vpc(),
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
			Architecture: config.LambdaConfig.Architecture,
			Environment:  &env,
			Tracing:      awslambda.Tracing_ACTIVE,
		},
	)

	database.Secret().GrantRead(function, nil)
	database.Connections().AllowFrom(
		function,
		awsec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow access from go lambda to database"),
	)

	return function
}

func NewVpc(stack awscdk.Stack) awsec2.Vpc {
	vpc := awsec2.NewVpc(
		stack, jsii.String("CustomVpc"), &awsec2.VpcProps{
			MaxAzs: jsii.Number(2),
		},
	)

	return vpc
}

func NewVpcEndpoint(stack awscdk.Stack, vpc awsec2.Vpc) awsec2.VpcEndpoint {
	return awsec2.NewInterfaceVpcEndpoint(
		stack, jsii.String("SecretManagerVpcEndpoint"), &awsec2.InterfaceVpcEndpointProps{
			Vpc:     vpc,
			Service: awsec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
		},
	)
}

func NewDatabase(stack awscdk.Stack, vpc awsec2.Vpc) awsrds.DatabaseInstance {
	database := awsrds.NewDatabaseInstance(
		stack, jsii.String("SmallDatabase"), &awsrds.DatabaseInstanceProps{
			Vpc:                     vpc,
			BackupRetention:         awscdk.Duration_Days(jsii.Number(0)),
			CloudwatchLogsRetention: awslogs.RetentionDays_ONE_MONTH,
			Engine: awsrds.DatabaseInstanceEngine_Postgres(
				&awsrds.PostgresInstanceEngineProps{
					Version: awsrds.PostgresEngineVersion_VER_16_1(),
				},
			),
			InstanceType:  awsec2.InstanceType_Of(awsec2.InstanceClass_T4G, awsec2.InstanceSize_MICRO),
			MultiAz:       jsii.Bool(false),
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		},
	)

	return database
}
