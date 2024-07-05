package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accessPolicyResource{}
	_ resource.ResourceWithConfigure   = &accessPolicyResource{}
	_ resource.ResourceWithImportState = &accessPolicyResource{}
)

// NewAccessPolicyResource is a helper function to simplify the provider implementation.
func NewAccessPolicyResource() resource.Resource {
	return &accessPolicyResource{}
}

// accessPolicyResource is the resource implementation.
type accessPolicyResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *accessPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_policy"
}

// Configure adds the provider configured client to the resource.
func (r *accessPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *accessPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Access Policy.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Access Policy.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Placeholder"),
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the Access Policy.",
				Optional:    true,
				Computed:    true,
			},
			"client_workload": schema.StringAttribute{
				Description: "Client workload ID configured in the Access Policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"trust_providers": schema.SetAttribute{
				Description: "Set of Trust Providers to enforce on the Access Policy.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"access_conditions": schema.SetAttribute{
				Description: "Set of Access Conditions to enforce on the Access Policy.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"credential_provider": schema.StringAttribute{
				Description: "Credential Provider ID configured in the Access Policy.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_workload": schema.StringAttribute{
				Description: "Server workload ID configured in the Access Policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *accessPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accessPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var policy aembit.PolicyDTO = convertAccessPolicyModelToPolicyDTO(plan, nil)

	// Create new Access Policy
	accessPolicy, err := r.client.CreateAccessPolicy(policy, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access policy",
			"Could not create access policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertAccessPolicyDTOToModel(*accessPolicy)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *accessPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state accessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed policy value from Aembit
	accessPolicy, err := r.client.GetAccessPolicy(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Access Policy",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state = convertAccessPolicyExternalDTOToModel(accessPolicy)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accessPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state accessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan accessPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var policy aembit.PolicyDTO = convertAccessPolicyModelToPolicyDTO(plan, &externalID)

	// Update Access Policy
	accessPolicy, err := r.client.UpdateAccessPolicy(policy, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access policy",
			"Could not update access policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertAccessPolicyDTOToModel(*accessPolicy)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accessPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state accessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Access Policy is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableAccessPolicy(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Access Policy",
				"Could not disable Access Policy, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Access Policy
	_, err := r.client.DeleteAccessPolicy(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Access Policy",
			"Could not delete access policy, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *accessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAccessPolicyModelToPolicyDTO(model accessPolicyResourceModel, externalID *string) aembit.PolicyDTO {
	var policy aembit.PolicyDTO
	policy.EntityDTO = aembit.EntityDTO{
		Name:     model.Name.ValueString(),
		IsActive: model.IsActive.ValueBool(),
	}
	policy.ClientWorkload = model.ClientWorkload.ValueString()
	policy.ServerWorkload = model.ServerWorkload.ValueString()
	policy.CredentialProvider = model.CredentialProvider.ValueString()

	if externalID != nil {
		policy.EntityDTO.ExternalID = *externalID
	}

	policy.TrustProviders = make([]string, len(model.TrustProviders))
	if model.TrustProviders != nil && len(model.TrustProviders) > 0 {
		for i, trustProvider := range model.TrustProviders {
			policy.TrustProviders[i] = trustProvider.ValueString()
		}
	}

	policy.AccessConditions = make([]string, len(model.AccessConditions))
	if model.AccessConditions != nil && len(model.AccessConditions) > 0 {
		for i, accessConditions := range model.AccessConditions {
			policy.AccessConditions[i] = accessConditions.ValueString()
		}
	}

	return policy
}

func convertAccessPolicyDTOToModel(dto aembit.PolicyDTO) accessPolicyResourceModel {
	var model accessPolicyResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload)

	if len(dto.CredentialProvider) > 0 {
		model.CredentialProvider = types.StringValue(dto.CredentialProvider)
	}

	model.TrustProviders = make([]types.String, len(dto.TrustProviders))
	if dto.TrustProviders != nil && len(dto.TrustProviders) > 0 {
		for i, trustProvider := range dto.TrustProviders {
			model.TrustProviders[i] = types.StringValue(trustProvider)
		}
	}

	model.AccessConditions = make([]types.String, len(dto.AccessConditions))
	if dto.AccessConditions != nil && len(dto.AccessConditions) > 0 {
		for i, accessConditions := range dto.AccessConditions {
			model.AccessConditions[i] = types.StringValue(accessConditions)
		}
	}

	return model
}

func convertAccessPolicyExternalDTOToModel(dto aembit.PolicyExternalDTO) accessPolicyResourceModel {
	var model accessPolicyResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload.ExternalID)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload.ExternalID)

	if len(dto.CredentialProvider.ExternalID) > 0 {
		model.CredentialProvider = types.StringValue(dto.CredentialProvider.ExternalID)
	}

	model.TrustProviders = make([]types.String, len(dto.TrustProviders))
	if dto.TrustProviders != nil && len(dto.TrustProviders) > 0 {
		for i, trustProvider := range dto.TrustProviders {
			model.TrustProviders[i] = types.StringValue(trustProvider.ExternalID)
		}
	}

	model.AccessConditions = make([]types.String, len(dto.AccessConditions))
	if dto.AccessConditions != nil && len(dto.AccessConditions) > 0 {
		for i, accessConditions := range dto.AccessConditions {
			model.AccessConditions[i] = types.StringValue(accessConditions.ExternalID)
		}
	}

	return model
}
