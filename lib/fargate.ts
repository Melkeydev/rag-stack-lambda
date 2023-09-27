import * as cdk from "aws-cdk-lib";
import { StackProps } from "aws-cdk-lib";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";
import { HostedZone, ARecord, RecordTarget } from "aws-cdk-lib/aws-route53";
import {
  Certificate,
  CertificateValidation,
} from "aws-cdk-lib/aws-certificatemanager";
import { ApplicationProtocol } from "aws-cdk-lib/aws-elasticloadbalancingv2";
import { LoadBalancerTarget } from "aws-cdk-lib/aws-route53-targets";

export interface FargateProps extends StackProps {
  domainName: string;
  hostedZoneName: string;
  hostedZoneId: string;
  aRecordName: string;
}

export class Backend extends Construct {
  public apiUrl: string;
  constructor(scope: Construct, id: string, props: FargateProps) {
    super(scope, id);

    // Define the DynamoDB table
    const table = new dynamodb.Table(this, "MyTable", {
      partitionKey: { name: "username", type: dynamodb.AttributeType.STRING },
      tableName: "user-table-name-ecs",
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const vpc = new ec2.Vpc(this, "RagFargateVPC", {
      maxAzs: 2,
    });

    const publicZone = HostedZone.fromHostedZoneAttributes(
      this,
      "HttpsFargateAlbPublicZone",
      {
        zoneName: props.hostedZoneName,
        hostedZoneId: props.hostedZoneId,
      }
    );

    const certificate = new Certificate(this, "HttpsFargateAlbCertificate", {
      domainName: props.domainName,
      validation: CertificateValidation.fromDns(publicZone),
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
          certificate: certificate,
          protocol: ApplicationProtocol.HTTPS,
          redirectHTTP: true,
          memoryLimitMiB: 512,
          publicLoadBalancer: true,
        }
      );

    new ARecord(this, "HttpsFargateAlbARecord", {
      zone: publicZone,
      recordName: props.aRecordName,
      target: RecordTarget.fromAlias(
        new LoadBalancerTarget(fargateService.loadBalancer)
      ),
    });

    // Grant ECS instance function
    table.grantReadWriteData(fargateService.taskDefinition.taskRole);
    this.apiUrl = props.aRecordName;
  }
}
