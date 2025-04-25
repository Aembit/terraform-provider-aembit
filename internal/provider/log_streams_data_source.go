package provider

import (
	"context"
	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
func (d *logStreamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *logStreamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_streams"
}

// Schema defines the schema for the resource.
func (d *logStreamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a log stream.",
		Attributes: map[string]schema.Attribute{
			"log_streams": schema.ListNestedAttribute{
				Description: "List of log streams.",
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
							Description: "Active/Inactive status of the Log Stream.",
							Computed:    true,
						},
						"data_type": schema.StringAttribute{
							Description: "Data type of Log Stream.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Destination type of Log Stream.",
							Required:    true,
						},
						"aws_s3_bucket": schema.SingleNestedAttribute{
							Description: "AWSS3Bucket destination type Log Stream configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"s3_bucket_region": schema.StringAttribute{
									Description: "S3 Bucket Region.",
									Required:    true,
								},
								"s3_bucket_name": schema.StringAttribute{
									Description: "S3 Bucket Name.",
									Required:    true,
								},
								"s3_path_prefix": schema.StringAttribute{
									Description: "S3 Path Prefix.",
									Required:    true,
								},
							},
						},
						"gcs_bucket": schema.SingleNestedAttribute{
							Description: "GCSBucket destination type Log Stream configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"gcs_bucket_name": schema.StringAttribute{
									Description: "GCS Bucket Name.",
									Required:    true,
								},
								"gcs_path_prefix": schema.StringAttribute{
									Description: "GCS Path Prefix.",
									Required:    true,
								},
								"audience": schema.StringAttribute{
									Description: "Audience.",
									Required:    true,
								},
								"service_account_email": schema.StringAttribute{
									Description: "Service Account Email.",
									Required:    true,
								},
								"token_lifetime": schema.Int64Attribute{
									Description: "Token Lifetime.",
									Required:    true,
								},
							},
						},
						"splunk_http_event_collector": schema.SingleNestedAttribute{
							Description: "SplunkHttpEventCollector destination type Log Stream configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"splunk_host_port": schema.StringAttribute{
									Description: "Splunk Host Port.",
									Required:    true,
								},
								"authentication_token": schema.StringAttribute{
									Description: "Authentication Token.",
									Required:    true,
								},
								"source_name": schema.StringAttribute{
									Description: "Source Name.",
									Required:    true,
								},
								"tls": schema.BoolAttribute{
									Description: "Tls.",
									Required:    true,
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
func (d *logStreamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
