import { Source, BucketDeployment } from "aws-cdk-lib/aws-s3-deployment";
import { Construct } from "constructs";
import { execSync, ExecSyncOptions } from "child_process";
import { RemovalPolicy, DockerImage } from "aws-cdk-lib";
import {
  Distribution,
  SecurityPolicyProtocol,
  ViewerProtocolPolicy,
  CachePolicy,
} from "aws-cdk-lib/aws-cloudfront";
import { S3Origin } from "aws-cdk-lib/aws-cloudfront-origins";
import { Bucket } from "aws-cdk-lib/aws-s3";
import * as fsExtra from "fs-extra";

interface FrontendProps {
  apiUrl: string;
}

export class Frontend extends Construct {
  public readonly siteBucket: Bucket;
  public readonly distribution: Distribution;

  constructor(scope: Construct, id: string, props: FrontendProps) {
    super(scope, id);

    this.siteBucket = new Bucket(this, "frontEndBucket", {
      publicReadAccess: false,
      removalPolicy: RemovalPolicy.DESTROY,
    });

    this.distribution = new Distribution(this, "CloudFrontDistribution", {
      enableLogging: true,
      minimumProtocolVersion: SecurityPolicyProtocol.TLS_V1_2_2021,
      defaultBehavior: {
        origin: new S3Origin(this.siteBucket),
        viewerProtocolPolicy: ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachePolicy: CachePolicy.CACHING_DISABLED,
      },
      defaultRootObject: "index.html",
    });

    const execOptions: ExecSyncOptions = { stdio: "inherit" };

    const bundle = Source.asset("./frontend", {
      bundling: {
        command: [
          "sh",
          "-c",
          'echo "Docker build not supported. Please install esbuild."',
        ],
        image: DockerImage.fromRegistry("alpine"),
        local: {
          /* istanbul ignore next */
          tryBundle(outputDir: string) {
            try {
              execSync("esbuild --version", execOptions);
            } catch {
              return false;
            }
            execSync(
              "cd frontend && npm install --legacy-peer-deps && npm run build",
              execOptions
            );
            fsExtra.copySync("./frontend/dist", outputDir, {
              ...execOptions,
              recursive: true,
            });
            return true;
          },
        },
      },
    });

    const config = {
      apiUrl: props.apiUrl,
    };

    new BucketDeployment(this, "DeployBucket", {
      sources: [bundle, Source.jsonData("config.json", config)],
      destinationBucket: this.siteBucket,
      distribution: this.distribution,
      distributionPaths: ["/*"],
    });
  }
}
