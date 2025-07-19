package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-aembit/internal/provider/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &logStreamsDataSource{}
	_ datasource.DataSourceWithConfigure = &logStreamsDataSource{}
)

// NewLogStreamsDataSource is a helper function to simplify the provider implementation.
func NewLogStreamsDataSource() datasource.DataSource {
	return &logStreamsDataSource{}
}

// logStreamsDataSource is the data source implementation.
type logStreamsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the logStream to the data source.
func (d *logStreamsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *logStreamsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_log_streams"
}

// Schema defines the schema for the resource.
func (d *logStreamsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about all Aembit Log Streams",
		Attributes: map[string]schema.Attribute{
			"log_streams": schema.ListNestedAttribute{
				Description: "List of Log Streams.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the Log Stream.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the Log Stream.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the Log Stream.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Status of the Log Stream (`active` or `inactive`)",
							Computed:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "Data type of Log Stream.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Destination type of Log Stream.",
							Computed:    true,
						},
						"aws_s3_bucket": schema.SingleNestedAttribute{
							Description: "Log Stream configuration for the AWS S3 Bucket destination type.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"s3_bucket_region": schema.StringAttribute{
									Description: "S3 Bucket Region.",
									Computed:    true,
								},
								"s3_bucket_name": schema.StringAttribute{
									Description: "S3 Bucket Name.",
									Computed:    true,
								},
								"s3_path_prefix": schema.StringAttribute{
									Description: "S3 Path Prefix.",
									Computed:    true,
								},
							},
						},
						"gcs_bucket": schema.SingleNestedAttribute{
							Description: "GCSBucket destination type Log Stream configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"gcs_bucket_name": schema.StringAttribute{
									Description: "GCS Bucket Name.",
									Computed:    true,
								},
								"gcs_path_prefix": schema.StringAttribute{
									Description: "GCS Path Prefix.",
									Computed:    true,
								},
								"audience": schema.StringAttribute{
									Description: "Audience.",
									Computed:    true,
								},
								"service_account_email": schema.StringAttribute{
									Description: "Service Account Email.",
									Computed:    true,
								},
								"token_lifetime": schema.Int64Attribute{
									Description: "Token Lifetime.",
									Computed:    true,
								},
							},
						},
						"splunk_http_event_collector": schema.SingleNestedAttribute{
							Description: "Log Stream configuration for the Splunk HTTP Event Collector (HEC) destination type.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"host_port": schema.StringAttribute{
									Description: "Splunk HTTP Event Collector (HEC) host:port value.",
									Computed:    true,
								},
								"source_name": schema.StringAttribute{
									Description: "Splunk Data Input Source Name.",
									Computed:    true,
								},
								"tls": schema.BoolAttribute{
									Description: "Splunk HTTP Event Collector (HEC) TLS configuration.",
									Computed:    true,
								},
								"tls_verification": schema.StringAttribute{
									Description: "Splunk HTTP Event Collector (HEC) TLS verification.",
									Computed:    true,
								},
							},
						},
						"crowdstrike_http_event_collector": schema.SingleNestedAttribute{
							Description: "Log Stream configuration for the Crowdstrike HTTP Event Collector (HEC) destination type.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"host_port": schema.StringAttribute{
									Description: "Crowdstrike HTTP Event Collector (HEC) host:port value.",
									Computed:    true,
								},
								"source_name": schema.StringAttribute{
									Description: "Crowdstrike Data Input Source Name.",
									Computed:    true,
								},
								"tls": schema.BoolAttribute{
									Description: "Crowdstrike HTTP Event Collector (HEC) TLS configuration.",
									Computed:    true,
								},
								"tls_verification": schema.StringAttribute{
									Description: "Crowdstrike HTTP Event Collector (HEC) TLS verification.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *logStreamsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.LogStreamsDataSourceModel

	logStreams, err := d.client.GetLogStreams(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Log Streams",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, logStream := range logStreams {
		logStreamState := convertLogStreamDTOToModel(logStream, models.LogStreamResourceModel{})
		state.LogStreams = append(state.LogStreams, logStreamState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
