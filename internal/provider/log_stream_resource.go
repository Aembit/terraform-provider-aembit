package provider

import (
	"context"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logStreamResource{}
	_ resource.ResourceWithConfigure   = &logStreamResource{}
	_ resource.ResourceWithImportState = &logStreamResource{}
)

// NewLogStreamResource is a helper function to simplify the log stream implementation.
func NewLogStreamResource() resource.Resource {
	return &logStreamResource{}
}

// logStreamResource is the resource implementation.
type logStreamResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *logStreamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_stream"
}

// Configure adds the log stream to the resource.
func (r *logStreamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *logStreamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Log Stream.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Log Stream.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Log Stream.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the Log Stream.",
				Optional:    true,
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
				Description: "Log Stream configuration for the AWS S3 Bucket destination type.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"s3_bucket_region": schema.StringAttribute{
						Description: "S3 Bucket Region.",
						Optional:    true,
						Validators: []validator.String{
							validators.S3BucketRegionValidation(),
							validators.S3BucketRegionLengthValidation(),
						},
					},
					"s3_bucket_name": schema.StringAttribute{
						Description: "S3 Bucket Name.",
						Optional:    true,
						Validators: []validator.String{
							validators.S3BucketNameValidation(),
							validators.S3BucketNameLengthValidation(),
						},
					},
					"s3_path_prefix": schema.StringAttribute{
						Description: "S3 Path Prefix.",
						Optional:    true,
						Validators: []validator.String{
							validators.S3PathPrefixValidation(),
							validators.S3PathPrefixLengthValidation(),
						},
					},
				},
			},
			"gcs_bucket": schema.SingleNestedAttribute{
				Description: "GCSBucket destination type Log Stream configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"gcs_bucket_name": schema.StringAttribute{
						Description: "GCS Bucket Name.",
						Optional:    true,
						Validators: []validator.String{
							validators.GCSBucketNameSimpleValidation(),
							validators.GCSBucketNameWithPeriodValidation(),
						},
					},
					"gcs_path_prefix": schema.StringAttribute{
						Description: "GCS Path Prefix.",
						Optional:    true,
						Validators: []validator.String{
							validators.GCSPathPrefixValidation(),
						},
					},
					"audience": schema.StringAttribute{
						Description: "Audience.",
						Optional:    true,
						Validators: []validator.String{
							validators.AudienceValidation(),
						},
					},
					"service_account_email": schema.StringAttribute{
						Description: "Service Account Email.",
						Optional:    true,
						Validators: []validator.String{
							validators.EmailValidation(),
						},
					},
					"token_lifetime": schema.Int64Attribute{
						Description: "Token Lifetime.",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.Between(300, 3600),
						},
					},
				},
			},
			"splunk_http_event_collector": schema.SingleNestedAttribute{
				Description: "Log Stream configuration for the Splunk HTTP Event Collector (HEC) destination type.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"hec_host_port": schema.StringAttribute{
						Description: "Splunk HTTP Event Collector (HEC) host:port value.",
						Required:    true,
						Validators: []validator.String{
							validators.HecHostPortValidation(),
						},
					},
					"authentication_token": schema.StringAttribute{
						Description: "Authentication token.",
						Sensitive:   true,
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							validators.AuthenticationTokenValidation(),
						},
					},
					"hec_source_name": schema.StringAttribute{
						Description: "Splunk Data Input Source Name.",
						Required:    true,
					},
					"tls": schema.BoolAttribute{
						Description: "Splunk HTTP Event Collector (HEC) TLS configuration.",
						Optional:    true,
						Default:     booldefault.StaticBool(true),
						Computed:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *logStreamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.LogStreamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.LogStreamDTO = convertLogStreamModelToDTO(plan, nil)

	// Create new log stream
	logStream, err := r.client.CreateLogStream(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Log Stream",
			"Could not create Log Stream, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertLogStreamDTOToModel(*logStream, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *logStreamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.LogStreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed log stream from Aembit
	logStream, err, notFound := r.client.GetLogStream(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Log Stream",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertLogStreamDTOToModel(logStream, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *logStreamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state models.LogStreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.LogStreamResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.LogStreamDTO = convertLogStreamModelToDTO(plan, &externalID)

	// Update Log Stream
	logStream, err := r.client.UpdateLogStream(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Log Stream",
			"Could not update Log Stream, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertLogStreamDTOToModel(*logStream, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *logStreamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.LogStreamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Log Stream
	_, err := r.client.DeleteLogStream(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Log Stream",
			"Could not delete Log Stream, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *logStreamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertLogStreamModelToDTO(model models.LogStreamResourceModel, externalID *string) aembit.LogStreamDTO {
	var logStream aembit.LogStreamDTO
	logStream.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	if externalID != nil {
		logStream.EntityDTO.ExternalID = *externalID
	}

	logStream.DataType = model.DataType.ValueString()
	logStream.Type = model.Type.ValueString()

	if model.AWSS3Bucket != nil {
		logStream.S3BucketRegion = model.AWSS3Bucket.S3BucketRegion.ValueString()
		logStream.S3BucketName = model.AWSS3Bucket.S3BucketName.ValueString()
		logStream.S3PathPrefix = model.AWSS3Bucket.S3PathPrefix.ValueString()
	}

	if model.GCSBucket != nil {
		logStream.GCSBucketName = model.GCSBucket.GCSBucketName.ValueString()
		logStream.GCSPathPrefix = model.GCSBucket.GCSPathPrefix.ValueString()
		logStream.Audience = model.GCSBucket.Audience.ValueString()
		logStream.ServiceAccountEmail = model.GCSBucket.ServiceAccountEmail.ValueString()
		logStream.TokenLifetime = model.GCSBucket.TokenLifetime.ValueInt64()
	}

	if model.SplunkHttpEventCollector != nil {
		logStream.HecHostPort = model.SplunkHttpEventCollector.HecHostPort.ValueString()
		logStream.AuthenticationToken = model.SplunkHttpEventCollector.AuthenticationToken.ValueString()
		logStream.HecSourceName = model.SplunkHttpEventCollector.HecSourceName.ValueString()
		logStream.Tls = model.SplunkHttpEventCollector.Tls.ValueBool()
	}

	return logStream
}

func convertLogStreamDTOToModel(dto aembit.LogStreamDTO, state models.LogStreamResourceModel) models.LogStreamResourceModel {
	var model models.LogStreamResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	model.DataType = types.StringValue(dto.DataType)
	model.Type = types.StringValue(dto.Type)

	if dto.Type == "AwsS3Bucket" {
		model.AWSS3Bucket = &models.AWSS3BucketModel{
			S3BucketRegion: types.StringValue(dto.S3BucketRegion),
			S3BucketName:   types.StringValue(dto.S3BucketName),
			S3PathPrefix:   types.StringValue(dto.S3PathPrefix),
		}
	}

	if dto.Type == "GcsBucket" {
		model.GCSBucket = &models.GCSBucketModel{
			GCSBucketName:       types.StringValue(dto.GCSBucketName),
			GCSPathPrefix:       types.StringValue(dto.GCSPathPrefix),
			Audience:            types.StringValue(dto.Audience),
			ServiceAccountEmail: types.StringValue(dto.ServiceAccountEmail),
			TokenLifetime:       types.Int64Value(dto.TokenLifetime),
		}
	}

	if dto.Type == "SplunkHttpEventCollector" {
		model.SplunkHttpEventCollector = &models.SplunkHttpEventCollectorModel{
			HecHostPort:   types.StringValue(dto.HecHostPort),
			HecSourceName: types.StringValue(dto.HecSourceName),
			Tls:           types.BoolValue(dto.Tls),
		}

		if dto.AuthenticationToken != "" {
			model.SplunkHttpEventCollector.AuthenticationToken = types.StringValue(dto.AuthenticationToken)
		} else if state.SplunkHttpEventCollector != nil {
			model.SplunkHttpEventCollector.AuthenticationToken = state.SplunkHttpEventCollector.AuthenticationToken
		} else {
			model.SplunkHttpEventCollector.AuthenticationToken = types.StringNull()
		}
	}

	return model
}
