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
func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration"
}

// Configure adds the provider configured client to the resource.
func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *integrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
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
					"audience": schema.StringAttribute{
						Description: "Audience for the OAuth Endpoint of the Integration.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.IntegrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.IntegrationDTO = convertIntegrationModelToDTO(ctx, plan, nil)

	// Create new Integration
	integration, err := r.client.CreateIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Integration",
			"Could not create Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertIntegrationDTOToModel(ctx, *integration, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.IntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	integration, err, notFound := r.client.GetIntegration(state.ID.ValueString(), nil)
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

	state = convertIntegrationDTOToModel(ctx, integration, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	var dto aembit.IntegrationDTO = convertIntegrationModelToDTO(ctx, plan, &externalID)

	// Update Integration
	integration, err := r.client.UpdateIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Integration",
			"Could not update Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertIntegrationDTOToModel(ctx, *integration, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.IntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Integration is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableIntegration(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Client Workload",
				"Could not disable Client Workload, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Integration
	_, err := r.client.DeleteIntegration(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Integration",
			"Could not delete Integration, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIntegrationModelToDTO(ctx context.Context, model models.IntegrationResourceModel, externalID *string) aembit.IntegrationDTO {
	var integration aembit.IntegrationDTO
	integration.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			integration.Tags = append(integration.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	if externalID != nil {
		integration.EntityDTO.ExternalID = *externalID
	}

	integration.Endpoint = model.Endpoint.ValueString()
	integration.Type = model.Type.ValueString()
	integration.SyncFrequencySeconds = model.SyncFrequency.ValueInt64()
	integration.IntegrationJSON = aembit.IntegrationJSONDTO{
		TokenURL:     model.OAuthClientCredentials.TokenURL.ValueString(),
		ClientID:     model.OAuthClientCredentials.ClientID.ValueString(),
		ClientSecret: model.OAuthClientCredentials.ClientSecret.ValueString(),
		Audience:     model.OAuthClientCredentials.Audience.ValueString(),
	}

	return integration
}

func convertIntegrationDTOToModel(ctx context.Context, dto aembit.IntegrationDTO, state models.IntegrationResourceModel) models.IntegrationResourceModel {
	var model models.IntegrationResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	model.Type = types.StringValue(dto.Type)
	model.Endpoint = types.StringValue(dto.Endpoint)
	model.SyncFrequency = types.Int64Value(dto.SyncFrequencySeconds)
	model.OAuthClientCredentials = &models.IntegrationOAuthClientCredentialsModel{
		TokenURL:     types.StringValue(dto.IntegrationJSON.TokenURL),
		ClientID:     types.StringValue(dto.IntegrationJSON.ClientID),
		ClientSecret: types.StringValue(dto.IntegrationJSON.ClientSecret),
		Audience:     types.StringValue(dto.IntegrationJSON.Audience),
	}
	if len(dto.IntegrationJSON.ClientSecret) == 0 && state.OAuthClientCredentials != nil {
		model.OAuthClientCredentials.ClientSecret = state.OAuthClientCredentials.ClientSecret
	}

	return model
}
