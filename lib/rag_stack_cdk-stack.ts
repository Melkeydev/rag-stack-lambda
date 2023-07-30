import * as cdk from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
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
    const api = new apigw.RestApi(this, "Endpoint");
    const integration = new apigw.LambdaIntegration(myFunction);
    api.root.addMethod("POST", integration);

    // Define the '/register' resource and method
    const registerResource = api.root.addResource("register");
    const registerIntegration = new apigw.LambdaIntegration(myFunction);
    registerResource.addMethod("POST", registerIntegration);

    const defaultIntegration = new apigw.LambdaIntegration(myFunction);
    api.root.addMethod("ANY", defaultIntegration, {
      methodResponses: [{ statusCode: "200" }, { statusCode: "404" }],
    });
  }
}
