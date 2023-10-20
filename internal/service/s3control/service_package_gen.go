// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package s3control

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	s3control_sdkv2 "github.com/aws/aws-sdk-go-v2/service/s3control"
	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	s3control_sdkv1 "github.com/aws/aws-sdk-go/service/s3control"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceAccessGrantsInstance,
			Name:    "Access Grants Instance",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "id",
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceAccountPublicAccessBlock,
			TypeName: "aws_s3_account_public_access_block",
		},
		{
			Factory:  dataSourceMultiRegionAccessPoint,
			TypeName: "aws_s3control_multi_region_access_point",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceAccessPoint,
			TypeName: "aws_s3_access_point",
		},
		{
			Factory:  resourceAccountPublicAccessBlock,
			TypeName: "aws_s3_account_public_access_block",
		},
		{
			Factory:  resourceAccessPointPolicy,
			TypeName: "aws_s3control_access_point_policy",
		},
		{
			Factory:  resourceBucket,
			TypeName: "aws_s3control_bucket",
			Name:     "Bucket",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  resourceBucketLifecycleConfiguration,
			TypeName: "aws_s3control_bucket_lifecycle_configuration",
		},
		{
			Factory:  resourceBucketPolicy,
			TypeName: "aws_s3control_bucket_policy",
		},
		{
			Factory:  resourceMultiRegionAccessPoint,
			TypeName: "aws_s3control_multi_region_access_point",
		},
		{
			Factory:  resourceMultiRegionAccessPointPolicy,
			TypeName: "aws_s3control_multi_region_access_point_policy",
		},
		{
			Factory:  resourceObjectLambdaAccessPoint,
			TypeName: "aws_s3control_object_lambda_access_point",
		},
		{
			Factory:  resourceObjectLambdaAccessPointPolicy,
			TypeName: "aws_s3control_object_lambda_access_point_policy",
		},
		{
			Factory:  resourceStorageLensConfiguration,
			TypeName: "aws_s3control_storage_lens_configuration",
			Name:     "Storage Lens Configuration",
			Tags:     &types.ServicePackageResourceTags{},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.S3Control
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context, config map[string]any) (*s3control_sdkv1.S3Control, error) {
	sess := config["session"].(*session_sdkv1.Session)

	return s3control_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(config["endpoint"].(string))})), nil
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*s3control_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return s3control_sdkv2.NewFromConfig(cfg, func(o *s3control_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
