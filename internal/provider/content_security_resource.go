package provider

import (
	"context"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &contentSecurityResource{}
	_ resource.ResourceWithConfigure      = &contentSecurityResource{}
	_ resource.ResourceWithImportState    = &contentSecurityResource{}
	_ resource.ResourceWithModifyPlan     = &contentSecurityResource{}
	_ resource.ResourceWithValidateConfig = &contentSecurityResource{}
)

const (
	errInvalidConfiguration = "Invalid Configuration"
)

// NewContentSecurityResource is a helper function to simplify the provider implementation.
func NewContentSecurityResource() resource.Resource {
	return &contentSecurityResource{}
}

// contentSecurityResource is the resource implementation.
type contentSecurityResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *contentSecurityResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_content_security"
}

// Configure adds the provider configured client to the resource.
func (r *contentSecurityResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *contentSecurityResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Content Security resource.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Content Security resource.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_set_id": schema.StringAttribute{
				Description: "ResourceSet unique identifier of the Content Security resource.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Content Security resource.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Content Security resource.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"is_active": schema.BoolAttribute{
				Description: "Status of the Content Security resource (`active` or `inactive`).",
				Required:    true,
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"type": schema.StringAttribute{
				Description: "Type of the Content Security resource. Currently supports: `CrowdStrikeAIDR` or `McpToolAccessControl`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("CrowdStrikeAIDR", "McpToolAccessControl"),
				},
			},
			"crowdstrike_falcon_aidr": schema.SingleNestedAttribute{
				Description: "CrowdStrike Falcon AIDR configuration settings.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"encrypted_token": schema.StringAttribute{
						Description: "The encrypted API token or client secret used to authenticate with the CrowdStrike Falcon AIDR service.",
						Required:    true,
						Sensitive:   true,
					},
					"base_url": schema.StringAttribute{
						Description: "The base URL of the CrowdStrike Falcon AIDR service endpoint or API gateway.",
						Required:    true,
						Validators: []validator.String{
							validators.SecureURLValidation(),
						},
					},
					"fail_open": schema.BoolAttribute{
						Description: "Indicates whether requests should be allowed (fail-open) or blocked (fail-closed) if the CrowdStrike Falcon AIDR service is unreachable.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"timeout_ms": schema.Int64Attribute{
						Description: "The connection timeout in milliseconds for requests sent to the CrowdStrike Falcon AIDR service.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(5000),
						Validators: []validator.Int64{
							int64validator.Between(5000, 600000),
						},
					},
					"max_retries": schema.Int64Attribute{
						Description: "The maximum number of retry attempts for requests to the CrowdStrike Falcon AIDR service.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(2),
						Validators: []validator.Int64{
							int64validator.Between(0, 10),
						},
					},
				},
			},
			"mcp_tool_access_control": schema.SingleNestedAttribute{
				Description: "MCP Tool Access Control configuration settings.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						Description: "The overall policy enforcement mode for MCP tool access control (e.g., Allow, Block).",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("Allow", "Block"),
						},
					},
					"visibility": schema.StringAttribute{
						Description: "The tool visibility control mode (e.g., AllowAll, AllowSpecific, BlockAll, BlockSpecific).",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("AllowAll", "AllowSpecific", "BlockAll", "BlockSpecific"),
						},
					},
					"invocation": schema.StringAttribute{
						Description: "The tool invocation control mode (e.g., AllowAll, AllowSpecific, BlockAll, BlockSpecific).",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("AllowAll", "AllowSpecific", "BlockAll", "BlockSpecific"),
						},
					},
				},
			},
		},
	}
}

func (r *contentSecurityResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
	modifyPlanForResourceSetId(ctx, req, resp, r.client)
}

func (r *contentSecurityResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("crowdstrike_falcon_aidr"),
			path.MatchRoot("mcp_tool_access_control"),
		),
	}
}

func (r *contentSecurityResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config models.ContentSecurityResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Type.IsUnknown() {
		return
	}

	csType := config.Type.ValueString()

	if csType == "CrowdStrikeAIDR" {
		if config.CrowdStrikeFalconAIDR == nil {
			resp.Diagnostics.AddError(
				errInvalidConfiguration,
				"The 'crowdstrike_falcon_aidr' block must be configured when 'type' is set to 'CrowdStrikeAIDR'.",
			)
		}
		if config.McpToolAccessControl != nil {
			resp.Diagnostics.AddError(
				errInvalidConfiguration,
				"The 'mcp_tool_access_control' block cannot be configured when 'type' is set to 'CrowdStrikeAIDR'.",
			)
		}
	}

	if csType == "McpToolAccessControl" {
		if config.McpToolAccessControl == nil {
			resp.Diagnostics.AddError(
				errInvalidConfiguration,
				"The 'mcp_tool_access_control' block must be configured when 'type' is set to 'McpToolAccessControl'.",
			)
		}
		if config.CrowdStrikeFalconAIDR != nil {
			resp.Diagnostics.AddError(
				errInvalidConfiguration,
				"The 'crowdstrike_falcon_aidr' block cannot be configured when 'type' is set to 'McpToolAccessControl'.",
			)
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *contentSecurityResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.ContentSecurityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto := convertContentSecurityModelToDTO(ctx, plan, nil, r.client)

	resourceSetId := getResourceSetId(plan.ResourceSetID, r.client)
	// Create new ContentSecurity
	contentSecurity, err := r.client.CreateContentSecurity(dto, nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Content Security resource",
			"Could not create Content Security resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertContentSecurityDTOToModel(ctx, *contentSecurity, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *contentSecurityResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.ContentSecurityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := getResourceSetId(state.ResourceSetID, r.client)

	// Get refreshed value from Aembit
	contentSecurity, err, notFound := r.client.GetContentSecurity(state.ID.ValueString(), nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Content Security resource",
			"Could not read Aembit Content Security resource from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertContentSecurityDTOToModel(ctx, contentSecurity, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *contentSecurityResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	errorUpdateMessage := "Error updating Content Security resource"

	// Get current state
	var state models.ContentSecurityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.ContentSecurityResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.ResourceSetID.Equal(plan.ResourceSetID) {
		resp.Diagnostics.AddError(
			errorUpdateMessage,
			"Changing the ResourceSet of the resource is not supported.",
		)
		return
	}

	// Generate API request body from plan
	dto := convertContentSecurityModelToDTO(ctx, plan, &externalID, r.client)

	// Update ContentSecurity
	contentSecurity, err := r.client.UpdateContentSecurity(dto, nil, &dto.ResourceSet)
	if err != nil {
		resp.Diagnostics.AddError(
			errorUpdateMessage,
			"Could not update Content Security resource, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertContentSecurityDTOToModel(ctx, *contentSecurity, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *contentSecurityResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.ContentSecurityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := getResourceSetId(state.ResourceSetID, r.client)

	// Check if Content Security is Active - if it is, disable it first
	if state.IsActive.ValueBool() {
		_, err := r.client.DisableContentSecurity(state.ID.ValueString(), nil, &resourceSetId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Content Security resource",
				"Could not disable Content Security resource, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Content Security
	_, err := r.client.DeleteContentSecurity(ctx, state.ID.ValueString(), nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Content Security resource",
			"Could not delete Content Security resource, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *contentSecurityResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertContentSecurityModelToDTO(
	ctx context.Context,
	model models.ContentSecurityResourceModel,
	externalID *string,
	client *aembit.CloudClient,
) aembit.ContentSecurityDTO {
	var dto aembit.ContentSecurityDTO
	dto.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	dto.ResourceSet = getResourceSetId(model.ResourceSetID, client)
	dto.Tags = collectAllTagsDto(ctx, client.DefaultTags, model.Tags)

	if externalID != nil {
		dto.ExternalID = *externalID
	}

	dto.Type = model.Type.ValueString()

	if model.CrowdStrikeFalconAIDR != nil {
		dto.EncryptedToken = model.CrowdStrikeFalconAIDR.EncryptedToken.ValueString()
		dto.BaseUrl = model.CrowdStrikeFalconAIDR.BaseUrl.ValueString()
		dto.FailOpen = model.CrowdStrikeFalconAIDR.FailOpen.ValueBool()
		dto.TimeoutMs = int(model.CrowdStrikeFalconAIDR.TimeoutMs.ValueInt64())
		dto.MaxRetries = int(model.CrowdStrikeFalconAIDR.MaxRetries.ValueInt64())
	}

	if model.McpToolAccessControl != nil {
		dto.Mode = model.McpToolAccessControl.Mode.ValueString()
		dto.Visibility = model.McpToolAccessControl.Visibility.ValueString()
		dto.Invocation = model.McpToolAccessControl.Invocation.ValueString()
	}

	return dto
}

func convertContentSecurityDTOToModel(
	ctx context.Context,
	dto aembit.ContentSecurityDTO,
	planModel *models.ContentSecurityResourceModel,
) models.ContentSecurityResourceModel {
	var model models.ContentSecurityResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.ResourceSetID = types.StringValue(dto.ResourceSet)
	model.Type = types.StringValue(dto.Type)

	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	if dto.Type == "CrowdStrikeAIDR" {
		model.CrowdStrikeFalconAIDR = &models.CrowdStrikeFalconAIDRContentSecurityModel{}
		model.CrowdStrikeFalconAIDR.BaseUrl = types.StringValue(dto.BaseUrl)
		model.CrowdStrikeFalconAIDR.FailOpen = types.BoolValue(dto.FailOpen)
		model.CrowdStrikeFalconAIDR.TimeoutMs = types.Int64Value(int64(dto.TimeoutMs))
		model.CrowdStrikeFalconAIDR.MaxRetries = types.Int64Value(int64(dto.MaxRetries))
		model.CrowdStrikeFalconAIDR.EncryptedToken = types.StringNull()

		if planModel.CrowdStrikeFalconAIDR != nil {
			model.CrowdStrikeFalconAIDR.EncryptedToken = planModel.CrowdStrikeFalconAIDR.EncryptedToken
		}
	}

	if dto.Type == "McpToolAccessControl" {
		model.McpToolAccessControl = &models.McpToolAccessControlContentSecurityModel{}
		model.McpToolAccessControl.Mode = types.StringValue(dto.Mode)
		model.McpToolAccessControl.Visibility = types.StringValue(dto.Visibility)
		model.McpToolAccessControl.Invocation = types.StringValue(dto.Invocation)
	}

	return model
}
