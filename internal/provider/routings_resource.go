package provider

import (
	"context"

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
	_ resource.Resource                = &routingResource{}
	_ resource.ResourceWithConfigure   = &routingResource{}
	_ resource.ResourceWithImportState = &routingResource{}
)

// NewroutingResource is a helper function to simplify the provider implementation.
func NewRoutingResource() resource.Resource {
	return &routingResource{}
}

// routingResource is the resource implementation.
type routingResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *routingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing"
}

// Configure adds the provider configured client to the resource.
func (r *routingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *routingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Routing.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Routing.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Routing.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Routing.",
				Optional:    true,
				Computed:    true,
			},
			"proxyUrl": schema.StringAttribute{
				Description: "Proxy URL for the Routing.",
				Optional:    true,
				Computed:    true,
			},
			"resourceSetId": schema.StringAttribute{
				Description: "ResourceSet Id for the Routing.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *routingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan routingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.RoutingDTO = convertRoutingModelToDTO(plan, nil)

	// Create new routing
	routing, err := r.client.CreateRouting(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Routing",
			"Could not create Routing, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertRoutingDTOToModel(*routing, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *routingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state routingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	routing, err, notFound := r.client.GetRouting(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Routing",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertRoutingDTOToModel(routing, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state routingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan routingResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.RoutingDTO = convertRoutingModelToDTO(plan, &externalID)

	// Update routing
	routing, err := r.client.UpdateRouting(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Routing",
			"Could not update Routing, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertRoutingDTOToModel(*routing, state)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state routingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Routing is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableRouting(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Routing",
				"Could not disable Routing, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing routing
	_, err := r.client.DeleteRouting(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting routing",
			"Could not delete Routing, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *routingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertRoutingModelToDTO(model routingResourceModel, externalID *string) aembit.RoutingDTO {
	var routing aembit.RoutingDTO
	routing.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	routing.ProxyUrl = model.ProxyUrl.ValueString()
	routing.ResourceSetId = model.ResourceSetId.ValueString()

	if externalID != nil {
		routing.EntityDTO.ExternalID = *externalID
	}

	return routing
}

func convertRoutingDTOToModel(dto aembit.RoutingDTO, _ routingResourceModel) routingResourceModel {
	var model routingResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ProxyUrl = types.StringValue(dto.ProxyUrl)
	model.ResourceSetId = types.StringValue(dto.ResourceSetId)

	return model
}

func convertListRoutingDTOToModel(dto aembit.ListRoutingDTO, _ routingResourceModel) routingResourceModel {
	var model routingResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.ProxyUrl = types.StringValue(dto.ProxyUrl)
	model.ResourceSetName = types.StringValue(dto.ResourceSetName)

	return model
}
