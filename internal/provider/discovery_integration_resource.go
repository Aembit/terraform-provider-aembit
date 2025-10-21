package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &discoveryIntegrationResource{}
	_ resource.ResourceWithConfigure   = &discoveryIntegrationResource{}
	_ resource.ResourceWithImportState = &discoveryIntegrationResource{}
)

// NewDiscoveryIntegrationResource is a helper function to simplify the provider implementation.
func NewDiscoveryIntegrationResource() resource.Resource {
	return &discoveryIntegrationResource{}
}

// discoveryIntegrationResource is the resource implementation.
type discoveryIntegrationResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *discoveryIntegrationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_discovery_integration"
}

// Configure adds the provider configured client to the resource.
func (r *discoveryIntegrationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *discoveryIntegrationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the discovery integration.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "User-defined name of the discovery integration.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "User-defined description of the discovery integration.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the discovery integration.",
				Optional:    true,
				Computed:    true,
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"type": schema.StringAttribute{
				Description: "Type of discovery integration. The only accepted value is `WizIntegrationApi`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"WizIntegrationApi",
					}...),
				},
			},
			"sync_frequency_seconds": schema.Int64Attribute{
				Description: "Frequency (in seconds) for synchronizing the discovery integration. Accepted range: 300-3600 seconds",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3600),
				Validators: []validator.Int64{
					int64validator.Between(300, 3600),
				},
			},
			"last_sync": schema.StringAttribute{
				Description: "ISO 8601-formatted last sync date of the discovery integration.",
				Computed:    true,
			},
			"last_sync_status": schema.StringAttribute{
				Description: "Status of the last sync of the discovery integration.",
				Computed:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "Endpoint that performs the discovery integration.",
				Required:    true,
			},
			"wiz_integration": schema.SingleNestedAttribute{
				Description: "Wiz-specific properties for the discovery integration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{
						Description: "Token URL for the Wiz Endpoint of the discovery integration.",
						Required:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "Client ID for the Wiz Endpoint of the discovery integration.",
						Required:    true,
						Sensitive:   true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Client Secret for the Wiz Endpoint of the discovery integration.",
						Required:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience for the Wiz Endpoint of the discovery integration.",
						Required:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *discoveryIntegrationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.DiscoveryIntegrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	discoveryIntegration := convertDiscoveryIntegrationModelToDTO(ctx, plan, nil)

	// Create new Discovery Integration
	discoveryIntegrationResponse, err := r.client.CreateDiscoveryIntegration(
		discoveryIntegration,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Discovery Integration",
			"Could not create Discovery Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertDiscoveryIntegrationDTOToModel(ctx, *discoveryIntegrationResponse, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *discoveryIntegrationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.DiscoveryIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workload value from Aembit
	discoveryIntegration, err, notFound := r.client.GetDiscoveryIntegration(
		state.ID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Discovery Integration",
			fmt.Sprintf(
				"An error occurred while attempting to fetch the Discovery Integration with ID '%s' from Terraform state: %v",
				state.ID.ValueString(),
				err.Error(),
			),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// Overwrite items with refreshed state
	state = convertDiscoveryIntegrationDTOToModel(ctx, discoveryIntegration, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *discoveryIntegrationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.DiscoveryIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.DiscoveryIntegrationResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	discoveryIntegration := convertDiscoveryIntegrationModelToDTO(ctx, plan, &externalID)

	// Update Discovery Integration
	serverWorkload, err := r.client.UpdateDiscoveryIntegration(discoveryIntegration, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Discovery Integration",
			"Could not update Discovery Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertDiscoveryIntegrationDTOToModel(ctx, *serverWorkload, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *discoveryIntegrationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.DiscoveryIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Discovery Integration is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableDiscoveryIntegration(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Discovery Integration",
				"Could not disable Discovery Integration, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Discovery Integration
	_, err := r.client.DeleteDiscoveryIntegration(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Discovery Integration",
			"Could not delete Discovery Integration, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalID.
func (r *discoveryIntegrationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertDiscoveryIntegrationModelToDTO(
	ctx context.Context,
	model models.DiscoveryIntegrationResourceModel,
	externalID *string,
) aembit.DiscoveryIntegrationDTO {
	var discoveryIntegration aembit.DiscoveryIntegrationDTO
	discoveryIntegration.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	discoveryIntegration.Type = model.Type.ValueString()
	discoveryIntegration.SyncFrequencySeconds = model.SyncFrequencySeconds.ValueInt64()
	discoveryIntegration.Endpoint = model.Endpoint.ValueString()

	if externalID != nil {
		discoveryIntegration.ExternalID = *externalID
	}

	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			discoveryIntegration.Tags = append(discoveryIntegration.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	if model.Wiz != nil {
		jsonDto := aembit.DiscoveryIntegrationJSONDTO{
			TokenUrl:     model.Wiz.TokenUrl.ValueString(),
			ClientId:     model.Wiz.ClientId.ValueString(),
			ClientSecret: model.Wiz.ClientSecret.ValueString(),
			Audience:     model.Wiz.Audience.ValueString(),
		}

		jsonBytes, _ := json.Marshal(jsonDto)
		discoveryIntegration.DiscoveryIntegrationJSON = string(jsonBytes)
	}

	return discoveryIntegration
}

func convertDiscoveryIntegrationDTOToModel(
	ctx context.Context,
	dto aembit.DiscoveryIntegrationDTO,
	state models.DiscoveryIntegrationResourceModel,
) models.DiscoveryIntegrationResourceModel {
	var model models.DiscoveryIntegrationResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Tags = newTagsModel(ctx, dto.Tags)
	model.Type = types.StringValue(dto.Type)
	model.SyncFrequencySeconds = types.Int64Value(dto.SyncFrequencySeconds)
	model.LastSync = types.StringValue(dto.LastSync)
	model.LastSyncStatus = types.StringValue(dto.LastSyncStatus)
	model.Endpoint = types.StringValue(dto.Endpoint)

	switch dto.Type {
	case "WizIntegrationApi":
		var wizDto aembit.DiscoveryIntegrationJSONDTO
		err := json.Unmarshal([]byte(dto.DiscoveryIntegrationJSON), &wizDto)
		if err != nil {
			fmt.Println("Failed to parse DiscoveryIntegrationJSON: ", err)
			return model
		}
		model.Wiz = &models.DiscoveryIntegrationWizModel{
			TokenUrl:     types.StringValue(wizDto.TokenUrl),
			ClientId:     types.StringValue(wizDto.ClientId),
			ClientSecret: types.StringNull(),
			Audience:     types.StringValue(wizDto.Audience),
		}

		if state.Wiz != nil {
			model.Wiz.ClientSecret = state.Wiz.ClientSecret
		}
	}

	return model
}
