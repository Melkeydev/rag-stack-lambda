import * as cdk from "aws-cdk-lib";
import { StackProps } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";

import { RestApi, LambdaIntegration } from "aws-cdk-lib/aws-apigateway";

export class Backend extends Construct {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    // Define the DynamoDB table
    const table = new dynamodb.Table(this, "MyTable2", {
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      tableName: "user-table-name",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // Define the Lambda function
    const myFunction = new lambda.Function(this, "MyFunction2", {
      code: lambda.Code.fromAsset("lambda"),
      handler: "main",
      runtime: lambda.Runtime.GO_1_X,
      environment: {
        // Rename to user table
        TABLE_NAME: table.tableName,
      },
    });

    // Grant Lambda function
    table.grantReadWriteData(myFunction);

    // Define the API Gateway
    const api = new RestApi(this, "Endpoint2", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      },
    });

    const integration = new LambdaIntegration(myFunction);
    api.root.addMethod("POST", integration);

    // Define the '/register' resource and method
    const registerResource = api.root.addResource("register");
    registerResource.addMethod("POST", integration);

    // Define the '/login' resource and method
    const loginResource = api.root.addResource("login");
    loginResource.addMethod("POST", integration);

    // Define the '/protected' resource and method
    const protectedResource = api.root.addResource("protected");
    protectedResource.addMethod("GET", integration);

    // Define the '/refresh' resource and method
    const refreshResource = api.root.addResource("refresh");
    refreshResource.addMethod("GET", integration);

    // Define the '/test' resource and method
    const testResource = api.root.addResource("test");
    testResource.addMethod("GET", integration);
  }
}
