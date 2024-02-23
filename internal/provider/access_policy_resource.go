package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.CloudClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.AembitClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *accessPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the access policy.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the access policy.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the access policy.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the access policy.",
				Optional:    true,
				Computed:    true,
			},
			"client_workload": schema.StringAttribute{
				Description: "Client workload ID configured in the access policy.",
				Required:    true,
			},
			"server_workload": schema.StringAttribute{
				Description: "Server workload ID configured in the access policy.",
				Required:    true,
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
	var policy aembit.PolicyDTO = convertAccessPolicyModelToPolicyDTO(ctx, plan, nil)

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
	plan = convertAccessPolicyDTOToModel(ctx, *accessPolicy)

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
		resp.Diagnostics.AddError(
			"Error Reading Aembit Access Policy",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state = convertAccessPolicyExternalDTOToModel(ctx, accessPolicy)

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
	var policy aembit.PolicyDTO = convertAccessPolicyModelToPolicyDTO(ctx, plan, &externalID)

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
	state = convertAccessPolicyDTOToModel(ctx, *accessPolicy)

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

	// Check if Access Policy is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Access Policy",
			"Access Policy is active and cannot be deleted. Please mark the policy as inactive first.",
		)
		return
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

// Imports an existing resource by passing externalId
func (r *accessPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAccessPolicyModelToPolicyDTO(ctx context.Context, model accessPolicyResourceModel, externalID *string) aembit.PolicyDTO {
	var policy aembit.PolicyDTO
	policy.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	policy.ClientWorkload = model.ClientWorkload.ValueString()
	policy.ServerWorkload = model.ServerWorkload.ValueString()

	if externalID != nil {
		policy.EntityDTO.ExternalID = *externalID
	}

	return policy
}

func convertAccessPolicyDTOToModel(ctx context.Context, dto aembit.PolicyDTO) accessPolicyResourceModel {
	var model accessPolicyResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload)

	return model
}

func convertAccessPolicyExternalDTOToModel(ctx context.Context, dto aembit.PolicyExternalDTO) accessPolicyResourceModel {
	var model accessPolicyResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload.ExternalID)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload.ExternalID)

	return model
}
