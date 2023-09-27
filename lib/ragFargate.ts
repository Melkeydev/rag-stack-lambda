import { CfnOutput } from "aws-cdk-lib";
import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import { Backend } from "./fargate";
import { Frontend } from "./frontend";

export class RagStackFargate extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const backend = new Backend(this, "Backend", {
      domainName: "*.devhouse.dev",
      aRecordName: "server.devhouse.dev",
      hostedZoneId: "Z00960303IO6O2SU42RW5",
      hostedZoneName: "devhouse.dev",
    });

    const frontend = new Frontend(this, "Frontend", { apiUrl: backend.apiUrl });
    new CfnOutput(this, "DistributionUrl", { value: frontend.distributionUrl });
  }
}
