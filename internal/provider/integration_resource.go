package provider

import (
	"context"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &integrationResource{}
	_ resource.ResourceWithConfigure   = &integrationResource{}
	_ resource.ResourceWithImportState = &integrationResource{}
	_ resource.ResourceWithModifyPlan  = &integrationResource{}
)

// NewIntegrationResource is a helper function to simplify the provider implementation.
func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

// integrationResource is the resource implementation.
type integrationResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *integrationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

// Configure adds the provider configured client to the resource.
func (r *integrationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *integrationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Integration.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Integration.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Integration.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Integration.",
				Optional:    true,
				Computed:    true,
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"type": schema.StringAttribute{
				Description: "Type of Aembit Integration. Possible values are: `WizIntegrationApi` or `CrowdStrike`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"WizIntegrationApi",
						"CrowdStrike",
					}...),
				},
			},
			"sync_frequency": schema.Int64Attribute{
				Description: "Frequency to be used for synchronizing the Integration.",
				Required:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint to be used for performing the Integration.",
				Required:    true,
				Validators: []validator.String{
					validators.UrlSchemeValidation(),
				},
			},
			"oauth_client_credentials": schema.SingleNestedAttribute{
				Description: "OAuth Client Credentials authentication information for the Integration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{
						Description: "Token URL for the OAuth Endpoint of the Integration.",
						Required:    true,
						Validators: []validator.String{
							validators.UrlSchemeValidation(),
						},
					},
					"client_id": schema.StringAttribute{
						Description: "Client ID for the OAuth Endpoint of the Integration.",
						Required:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Client Secret for the OAuth Endpoint of the Integration.",
						Required:    true,
						Sensitive:   true,
					},
					"wiz_integration": schema.SingleNestedAttribute{
						Description: "Wiz integration configuration.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"audience": schema.StringAttribute{
								Description: "Audience for the Wiz Integration.",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (r *integrationResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
}

// Create creates the resource and sets the initial Terraform state.
func (r *integrationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.IntegrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto := convertIntegrationModelToDTO(ctx, plan, nil, r.client.DefaultTags)

	// Create new Integration
	integration, err := r.client.CreateIntegrationV2(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Integration",
			"Could not create Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertIntegrationDTOToModel(ctx, *integration, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *integrationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.IntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	integration, err, notFound := r.client.GetIntegrationV2(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Integration",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertIntegrationDTOToModel(ctx, integration, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *integrationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.IntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.IntegrationResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto := convertIntegrationModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)

	// Update Integration
	integration, err := r.client.UpdateIntegrationV2(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Integration",
			"Could not update Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertIntegrationDTOToModel(ctx, *integration, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *integrationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.IntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Integration is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableIntegrationV2(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Client Workload",
				"Could not disable Client Workload, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Integration
	_, err := r.client.DeleteIntegrationV2(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Integration",
			"Could not delete Integration, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *integrationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIntegrationModelToDTO(
	ctx context.Context,
	model models.IntegrationResourceModel,
	externalID *string,
	defaultTags map[string]string,
) aembit.IntegrationV2DTO {
	var integration aembit.IntegrationV2DTO
	integration.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	if externalID != nil {
		integration.ExternalID = *externalID
	}

	integration.Endpoint = model.Endpoint.ValueString()
	integration.Type = model.Type.ValueString()
	integration.IntegrationType = model.Type.ValueString()
	integration.SyncFrequencySeconds = model.SyncFrequency.ValueInt64()
	integration.TokenURL = model.OAuthClientCredentials.TokenURL.ValueString()
	integration.ClientID = model.OAuthClientCredentials.ClientID.ValueString()
	integration.ClientSecret = model.OAuthClientCredentials.ClientSecret.ValueString()

	if model.Type.ValueString() == "WizIntegrationApi" {
		if model.OAuthClientCredentials.WizIntegration != nil {
			if !model.OAuthClientCredentials.WizIntegration.Audience.IsNull() {
				integration.Audience = model.OAuthClientCredentials.WizIntegration.Audience.ValueString()
			}
		}
	}

	integration.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)
	return integration
}

func convertIntegrationDTOToModel(
	ctx context.Context,
	dto aembit.IntegrationV2DTO,
	planModel *models.IntegrationResourceModel,
) models.IntegrationResourceModel {
	var model models.IntegrationResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)

	model.Type = types.StringValue(dto.Type)
	model.Endpoint = types.StringValue(dto.Endpoint)
	model.SyncFrequency = types.Int64Value(dto.SyncFrequencySeconds)

	oauthModel := &models.IntegrationOAuthClientCredentialsModel{
		TokenURL:       types.StringValue(dto.TokenURL),
		ClientID:       types.StringValue(dto.ClientID),
		ClientSecret:   types.StringNull(),
		WizIntegration: nil,
	}

	if planModel.OAuthClientCredentials != nil && !planModel.OAuthClientCredentials.ClientSecret.IsNull() {
		oauthModel.ClientSecret = planModel.OAuthClientCredentials.ClientSecret
	}

	model.OAuthClientCredentials = oauthModel

	if dto.Type == "WizIntegrationApi" && dto.Audience != "" {
		model.OAuthClientCredentials.WizIntegration = &models.WizIntegrationModel{
			Audience: types.StringValue(dto.Audience),
		}
	}

	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)
	return model
}
