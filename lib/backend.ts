import * as cdk from "aws-cdk-lib";
import { StackProps } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";

import {
  RestApi,
  LambdaIntegration,
  MethodLoggingLevel,
  EndpointType,
} from "aws-cdk-lib/aws-apigateway";

export class Backend extends Construct {
  public apiUrl: string;
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    // Define the DynamoDB table
    const table = new dynamodb.Table(this, "MyTable", {
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      tableName: "user-table-name",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    // Define the Lambda function
    const myFunction = new lambda.Function(this, "MyFunction", {
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
    const api = new RestApi(this, "exampleAPI", {
      defaultCorsPreflightOptions: {
        allowHeaders: [
          "Content-Type",
          "X-Amz-Date",
          "Authorization",
          "X-Api-Key",
        ],
        allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
        allowCredentials: true,
        allowOrigins: ["*"],
      },
      deployOptions: {
        loggingLevel: MethodLoggingLevel.INFO,
        dataTraceEnabled: true,
      },
      endpointConfiguration: {
        types: [EndpointType.REGIONAL],
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

    this.apiUrl = api.url;
  }
}
