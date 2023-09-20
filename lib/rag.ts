import { CfnOutput } from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Backend } from "./backend";
import { Frontend } from "./frontend";

export class RagStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const backend = new Backend(this, "Backend");
    const frontend = new Frontend(this, "Frontend", { apiUrl: backend.apiUrl });

    new CfnOutput(this, "DistributionUrl", { value: frontend.distributionUrl });
  }
}
