import * as cdk from "aws-cdk-lib";
import { StackProps } from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";

export class Backend extends Construct {
  public apiUrl: string;
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    // Define the DynamoDB table
    const table = new dynamodb.Table(this, "MyTable", {
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      tableName: "user-table-name-ecs",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const vpc = new ec2.Vpc(this, "RagFargateVPC", {
      // I think 2 is cheapr but i literally made that up
      maxAzs: 2,
    });

    const cluster = new ecs.Cluster(this, "RagFargateCluster", {
      vpc: vpc,
    });

    // Create a load-balanced Fargate service and make it public
    const fargateService =
      new ecs_patterns.ApplicationLoadBalancedFargateService(
        this,
        "MyFargateService",
        {
          cluster: cluster,
          cpu: 256,
          desiredCount: 1,
          taskImageOptions: {
            image: ecs.ContainerImage.fromAsset("fargate"),
            containerPort: 8080,
            environment: {
              TABLE_NAME: table.tableName,
            },
          },
          memoryLimitMiB: 512,
          publicLoadBalancer: true,
        }
      );

    // Grant ECS instance function
    table.grantReadWriteData(fargateService.taskDefinition.taskRole);

    this.apiUrl = fargateService.loadBalancer.loadBalancerDnsName;
  }
}
