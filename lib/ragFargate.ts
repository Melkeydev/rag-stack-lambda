import { CfnOutput } from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Backend } from "./fargate";
import { Frontend } from "./frontend";

export class RagStackFargate extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const backend = new Backend(this, "Backend");
    const frontend = new Frontend(this, "Frontend", { apiUrl: "devhouse.dev" });

    new CfnOutput(this, "DistributionUrl", { value: frontend.distributionUrl });
  }
}
