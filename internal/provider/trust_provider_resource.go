package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &trustProviderResource{}
	_ resource.ResourceWithConfigure   = &trustProviderResource{}
	_ resource.ResourceWithImportState = &trustProviderResource{}
)

// NewTrustProviderResource is a helper function to simplify the provider implementation.
func NewTrustProviderResource() resource.Resource {
	return &trustProviderResource{}
}

// trustProviderResource is the resource implementation.
type trustProviderResource struct {
	client *aembit.Client
}

// Metadata returns the resource type name.
func (r *trustProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_provider"
}

// Configure adds the provider configured client to the resource.
func (r *trustProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *trustProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the trust provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the trust provider.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the trust provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the trust provider.",
				Optional:    true,
				Computed:    true,
			},
			"azure_metadata": schema.SingleNestedAttribute{
				Description: "Azure Metadata type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sku": schema.StringAttribute{
						Optional: true,
					},
					"vm_id": schema.StringAttribute{
						Optional: true,
					},
					"subscription_id": schema.StringAttribute{
						Optional: true,
						//Validators: []validator.String{
						//	// Validate azure_metadata has at least one value
						//	stringvalidator.AtLeastOneOf(path.Expressions{
						//		path.MatchRelative().AtParent().AtName("sku"),
						//		path.MatchRelative().AtParent().AtName("vm_id"),
						//	}...),
						//},
					},
				},
				Validators: []validator.Object{
					// Validate azure_metadata has at least one value
					objectvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("azure_metadata").AtName("sku"),
						path.MatchRoot("azure_metadata").AtName("vm_id"),
						path.MatchRoot("azure_metadata").AtName("subscription_id"),
					}...),
				},
			},
			//"azure_metadata": schema.ObjectAttribute{
			//	Optional:       true,
			//	AttributeTypes: trustProviderAzureMetadataModel.AttrTypes,
			//	Validators:     trustProviderAzureMetadataValidators,
			//},
		},
	}
}

// Configure validators to ensure that only one trust provider type is specified
func (r *trustProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("azure_metadata"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *trustProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, nil)

	// Create new Trust Provider
	trust_provider, err := r.client.CreateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust provider",
			"Could not create trust provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderDTOToModel(ctx, *trust_provider)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trustProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	trust_provider, err := r.client.GetTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Trust Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertTrustProviderDTOToModel(ctx, trust_provider)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	external_id := state.ID.ValueString()

	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, &external_id)

	// Update Trust Provider
	trust_provider, err := r.client.UpdateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating trust provider",
			"Could not update trust provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertTrustProviderDTOToModel(ctx, *trust_provider)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *trustProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Trust Provider is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider",
			"Trust Provider is active and cannot be deleted. Please mark the trust as inactive first.",
		)
		return
	}

	// Delete existing Trust Provider
	_, err := r.client.DeleteTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider",
			"Could not delete trust provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId
func (r *trustProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertTrustProviderModelToDTO(ctx context.Context, model trustProviderResourceModel, external_id *string) aembit.TrustProviderDTO {
	var trust aembit.TrustProviderDTO
	trust.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if external_id != nil {
		trust.EntityDTO.ExternalId = *external_id
	}

	// Handle the Azure Metadata use case
	if model.AzureMetadata != nil {
		trust.Provider = "AzureMetadataService"
		trust.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)

		if len(model.AzureMetadata.Sku.ValueString()) > 0 {
			trust.MatchRules = append(trust.MatchRules, aembit.TrustProviderMatchRuleDTO{
				Attribute: "AzureSku",
				Value:     model.AzureMetadata.Sku.ValueString(),
			})
		}
		if len(model.AzureMetadata.VmId.ValueString()) > 0 {
			trust.MatchRules = append(trust.MatchRules, aembit.TrustProviderMatchRuleDTO{
				Attribute: "AzureVmId",
				Value:     model.AzureMetadata.VmId.ValueString(),
			})
		}
		if len(model.AzureMetadata.SubscriptionId.ValueString()) > 0 {
			trust.MatchRules = append(trust.MatchRules, aembit.TrustProviderMatchRuleDTO{
				Attribute: "AzureSubscriptionId",
				Value:     model.AzureMetadata.SubscriptionId.ValueString(),
			})
		}
	}

	return trust
}

func convertTrustProviderDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO) trustProviderResourceModel {
	var model trustProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalId)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	switch dto.Provider {
	case "AzureMetadataService": // Azure Metadata
		model.AzureMetadata = &trustProviderAzureMetadataModel{
			Sku:            types.StringNull(),
			VmId:           types.StringNull(),
			SubscriptionId: types.StringNull(),
		}

		for _, rule := range dto.MatchRules {
			switch rule.Attribute {
			case "AzureSku":
				model.AzureMetadata.Sku = types.StringValue(rule.Value)
			case "AzureVmId":
				model.AzureMetadata.VmId = types.StringValue(rule.Value)
			case "AzureSubscriptionId":
				model.AzureMetadata.SubscriptionId = types.StringValue(rule.Value)
			}
		}
	}

	return model
}
