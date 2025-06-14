package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LogStreamResourceModel maps the resource schema.
type LogStreamResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                            types.String                        `tfsdk:"id"`
	Name                          types.String                        `tfsdk:"name"`
	Description                   types.String                        `tfsdk:"description"`
	IsActive                      types.Bool                          `tfsdk:"is_active"`
	DataType                      types.String                        `tfsdk:"data_type"`
	Type                          types.String                        `tfsdk:"type"`
	AWSS3Bucket                   *AWSS3BucketModel                   `tfsdk:"aws_s3_bucket"`
	GCSBucket                     *GCSBucketModel                     `tfsdk:"gcs_bucket"`
	SplunkHttpEventCollector      *SplunkHttpEventCollectorModel      `tfsdk:"splunk_http_event_collector"`
	CrowdstrikeHttpEventCollector *CrowdstrikeHttpEventCollectorModel `tfsdk:"crowdstrike_http_event_collector"`
}

// LogStreamsDataSourceModel maps the datasource schema.
type LogStreamsDataSourceModel struct {
	LogStreams []LogStreamResourceModel `tfsdk:"log_streams"`
}

type AWSS3BucketModel struct {
	S3BucketRegion types.String `tfsdk:"s3_bucket_region"`
	S3BucketName   types.String `tfsdk:"s3_bucket_name"`
	S3PathPrefix   types.String `tfsdk:"s3_path_prefix"`
}

type GCSBucketModel struct {
	GCSBucketName       types.String `tfsdk:"gcs_bucket_name"`
	GCSPathPrefix       types.String `tfsdk:"gcs_path_prefix"`
	Audience            types.String `tfsdk:"audience"`
	ServiceAccountEmail types.String `tfsdk:"service_account_email"`
	TokenLifetime       types.Int64  `tfsdk:"token_lifetime"`
}

type SplunkHttpEventCollectorModel struct {
	HecHostPort         types.String `tfsdk:"host_port"`
	AuthenticationToken types.String `tfsdk:"authentication_token"`
	HecSourceName       types.String `tfsdk:"source_name"`
	Tls                 types.Bool   `tfsdk:"tls"`
	TlsVerification     types.String `tfsdk:"tls_verification"`
}

type CrowdstrikeHttpEventCollectorModel struct {
	HecHostPort     types.String `tfsdk:"host_port"`
	APIKey          types.String `tfsdk:"api_key"`
	HecSourceName   types.String `tfsdk:"source_name"`
	Tls             types.Bool   `tfsdk:"tls"`
	TlsVerification types.String `tfsdk:"tls_verification"`
}
