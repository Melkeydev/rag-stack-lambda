import * as cdk from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import { Construct } from "constructs";

export class RagStackCdkStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Define the Lambda function
    const myFunction = new lambda.Function(this, "MyFunction", {
      code: lambda.Code.fromAsset("lambda"),
      handler: "main",
      runtime: lambda.Runtime.GO_1_X,
    });

    // Define the API Gateway
    const api = new apigw.RestApi(this, "Endpoint");
    const integration = new apigw.LambdaIntegration(myFunction);
    api.root.addMethod("POST", integration);
  }
}
