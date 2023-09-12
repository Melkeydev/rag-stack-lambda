import * as cdk from "aws-cdk-lib";
import { StackProps } from "aws-cdk-lib";
import * as s3 from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";

export class Frontend extends Construct {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id);

    const bucket = new s3.Bucket(this, "FEBucket", {
      publicReadAccess: false,
    });
  }
}
