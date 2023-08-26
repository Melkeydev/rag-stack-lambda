import * as cdk from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { VPC } from "./vpc";
import { RDS } from "./rds";
import { Construct } from "constructs";

export class RagStackCdkStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Define the DynamoDB table
    const table = new dynamodb.Table(this, "MyTable", {
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      tableName: "user-table-name",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // Define VPC
    const vpc = new VPC(this, "VPC");

    // Define the RDS Instance
    const rds = new RDS(this, "rds", {
      vpc: vpc.vpc,
      securityGroup: vpc.securityGroup,
    });

    // Define the Lambda function
    const myFunction = new lambda.Function(this, "MyFunction", {
      code: lambda.Code.fromAsset("lambda"),
      handler: "main",
      runtime: lambda.Runtime.GO_1_X,
      vpc: vpc.vpc,
      securityGroups: [vpc.securityGroup],
      environment: {
        // Rename to user table
        TABLE_NAME: table.tableName,
      },
    });

    // Grant Lambda function
    table.grantReadWriteData(myFunction);

    // Define the API Gateway
    const api = new apigw.RestApi(this, "Endpoint");

    const integration = new apigw.LambdaIntegration(myFunction);
    api.root.addMethod("POST", integration);

    // Define the '/register' resource and method
    const registerIntegration = new apigw.LambdaIntegration(myFunction);
    const registerResource = api.root.addResource("register");
    registerResource.addMethod("POST", registerIntegration);

    // Define the '/login' resource and method
    const loginIntegration = new apigw.LambdaIntegration(myFunction);
    const loginResource = api.root.addResource("login");
    loginResource.addMethod("POST", loginIntegration);

    // Define the '/protected' resource and method
    const protectedIntegration = new apigw.LambdaIntegration(myFunction);
    const protectedResource = api.root.addResource("protected");
    protectedResource.addMethod("GET", protectedIntegration);

    // Define the '/refresh' resource and method
    const refreshIntegration = new apigw.LambdaIntegration(myFunction);
    const refreshResource = api.root.addResource("refresh");
    refreshResource.addMethod("GET", refreshIntegration);

    // Define the '/seed' resource and method
    const seedIntegration = new apigw.LambdaIntegration(myFunction);
    const seedResource = api.root.addResource("seed");
    seedResource.addMethod("GET", seedIntegration);
  }
}
