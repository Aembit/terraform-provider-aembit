package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resourceSetResource{}
	_ resource.ResourceWithConfigure   = &resourceSetResource{}
	_ resource.ResourceWithImportState = &resourceSetResource{}
)

// NewResourceSetResource is a helper function to simplify the provider implementation.
func NewResourceSetResource() resource.Resource {
	return &resourceSetResource{}
}

// resourceSetResource is the resource implementation.
type resourceSetResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *resourceSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_set"
}

// Configure adds the provider configured client to the resource.
func (r *resourceSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *resourceSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the ResourceSet.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the ResourceSet.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the ResourceSet.",
				Optional:    true,
				Computed:    true,
			},
			"roles": schema.SetAttribute{
				Description: "IDs of the Roles associated with this ResourceSet.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resourceSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.ResourceSetDTO = convertResourceSetModelToDTO(ctx, plan, nil)

	// Create new ResourceSet
	role, err := r.client.CreateResourceSet(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ResourceSet",
			"Could not create ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertResourceSetDTOToModel(ctx, *role)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state resourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	role, err, notFound := r.client.GetResourceSet(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit ResourceSet",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertResourceSetDTOToModel(ctx, role)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state resourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan resourceSetResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.ResourceSetDTO = convertResourceSetModelToDTO(ctx, plan, &externalID)

	// Update ResourceSet
	role, err := r.client.UpdateResourceSet(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ResourceSet",
			"Could not update ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertResourceSetDTOToModel(ctx, *role)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state resourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing ResourceSet
	_, err := r.client.DeleteResourceSet(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ResourceSet",
			"Could not delete ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *resourceSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Model to DTO conversion methods.
func convertResourceSetModelToDTO(_ context.Context, model resourceSetResourceModel, externalID *string) aembit.ResourceSetDTO {
	var dto aembit.ResourceSetDTO
	dto.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    true,
	}
	if externalID != nil {
		dto.EntityDTO.ExternalID = *externalID
	}

	dto.Roles = make([]string, len(model.Roles))
	if model.Roles != nil && len(model.Roles) > 0 {
		for i, role := range model.Roles {
			dto.Roles[i] = role.ValueString()
		}
	}

	return dto
}

// DTO to Model conversion methods.
func convertResourceSetDTOToModel(_ context.Context, dto aembit.ResourceSetDTO) resourceSetResourceModel {
	var model resourceSetResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)

	model.Roles = make([]types.String, len(dto.Roles))
	if dto.Roles != nil && len(dto.Roles) > 0 {
		for i, role := range dto.Roles {
			model.Roles[i] = types.StringValue(role)
		}
	}

	return model
}
