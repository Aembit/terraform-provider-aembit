package provider

import (
	"context"
	"encoding/json"
	"slices"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const resourceSetImplicitRolesKey = "implicit_roles"

// privateStateData represents the provider private state reader and writer interface.
type privateStateData interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

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
func (r *resourceSetResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_resource_set"
}

// Configure adds the provider configured client to the resource.
func (r *resourceSetResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *resourceSetResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage Resource Sets in an Aembit Cloud tenant.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the ResourceSet.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the ResourceSet.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the ResourceSet.",
				Optional:    true,
				Computed:    true,
			},
			"roles": schema.SetAttribute{
				Description: "IDs of the Roles associated with this Resource Set. Aembit automatically assigns the tenant's SuperAdmin and Auditor roles to every Resource Set. If omitted or null, Terraform populates this attribute with those default system role IDs. If set to an explicit set of custom role IDs or an empty set ([]), Terraform tracks only the configured role IDs in this attribute while exposing the full backend-assigned set (including SuperAdmin and Auditor) in effective_roles.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.UUIDRegexValidation()),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"effective_roles": schema.SetAttribute{
				Description: "Complete set of Role IDs in effect for this Resource Set, including system-assigned roles (SuperAdmin and Auditor) and any configured custom roles.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"standalone_certificate_authority": schema.StringAttribute{
				Description: "Standalone Certificate Authority ID configured for this ResourceSet.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *resourceSetResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.ResourceSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planRoles := plan.Roles

	// Generate API request body from plan
	resourceSetDTO, dtoDiags := convertResourceSetModelToDTO(ctx, plan, nil)
	resp.Diagnostics.Append(dtoDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new ResourceSet
	resourceSet, err := r.client.CreateResourceSet(resourceSetDTO, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ResourceSet",
			"Could not create ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	var modelDiags diag.Diagnostics
	plan, modelDiags = convertResourceSetDTOToModel(ctx, *resourceSet)
	resp.Diagnostics.Append(modelDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reconcile planned roles with backend-injected roles and persist implicit roles in private state
	reconciledRoles, reconcileDiags := reconcileRolesOnApply(ctx, planRoles, resourceSet.Roles, resp.Private)
	resp.Diagnostics.Append(reconcileDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Roles = reconciledRoles

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *resourceSetResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.ResourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed ResourceSet value from Aembit
	resourceSet, err, notFound := r.client.GetResourceSet(state.ID.ValueString(), nil)
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

	var modelDiags diag.Diagnostics
	state, modelDiags = convertResourceSetDTOToModel(ctx, resourceSet)
	resp.Diagnostics.Append(modelDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Filter out implicit roles using private state data
	filteredRoles, filterDiags := filterRolesOnRead(ctx, resourceSet.Roles, req.Private)
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Roles = filteredRoles

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *resourceSetResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.ResourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.ResourceSetResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	resourceSetDTO, dtoDiags := convertResourceSetModelToDTO(ctx, plan, &externalID)
	resp.Diagnostics.Append(dtoDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update ResourceSet
	resourceSet, err := r.client.UpdateResourceSet(resourceSetDTO, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ResourceSet",
			"Could not update ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	var modelDiags diag.Diagnostics
	state, modelDiags = convertResourceSetDTOToModel(ctx, *resourceSet)
	resp.Diagnostics.Append(modelDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Reconcile planned roles with backend-injected roles and persist implicit roles in private state
	reconciledRoles, reconcileDiags := reconcileRolesOnApply(ctx, plan.Roles, resourceSet.Roles, resp.Private)
	resp.Diagnostics.Append(reconcileDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Roles = reconciledRoles

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *resourceSetResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.ResourceSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing ResourceSet
	_, err := r.client.DeleteResourceSet(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ResourceSet",
			"Could not delete ResourceSet, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *resourceSetResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Model to DTO conversion methods.
func convertResourceSetModelToDTO(
	ctx context.Context,
	model models.ResourceSetResourceModel,
	externalID *string,
) (aembit.ResourceSetDTO, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dto aembit.ResourceSetDTO
	dto.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    true,
	}
	if externalID != nil {
		dto.ExternalID = *externalID
	}

	if !model.Roles.IsNull() && !model.Roles.IsUnknown() {
		diags.Append(model.Roles.ElementsAs(ctx, &dto.Roles, false)...)
	}

	dto.StandaloneCertificateAuthority = model.StandaloneCertificateAuthority.ValueString()

	return dto, diags
}

// DTO to Model conversion methods.
func convertResourceSetDTOToModel(
	ctx context.Context,
	dto aembit.ResourceSetDTO,
) (models.ResourceSetResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var model models.ResourceSetResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)

	if dto.StandaloneCertificateAuthority == "" {
		model.StandaloneCertificateAuthority = types.StringNull()
	} else {
		model.StandaloneCertificateAuthority = types.StringValue(dto.StandaloneCertificateAuthority)
	}

	rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, dto.Roles)
	diags.Append(setDiags...)
	model.Roles = rolesSet

	effectiveRolesSet, effectiveSetDiags := types.SetValueFrom(ctx, types.StringType, dto.Roles)
	diags.Append(effectiveSetDiags...)
	model.EffectiveRoles = effectiveRolesSet

	return model, diags
}

// reconcileRolesOnApply reconciles planned roles with backend-returned roles during Create and Update.
// If planRoles is unknown or null (omitted or set to null), it accepts all backendRoles and marks implicit_roles as an empty JSON array.
// If planRoles is known and non-null, it computes the implicit roles injected by the backend (backendRoles \ plannedRoles),
// persists them to the provider private state, and returns the reconciled roles (backendRoles \ implicitRoles).
func reconcileRolesOnApply(
	ctx context.Context,
	planRoles types.Set,
	backendRoles []string,
	privateState privateStateData,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planRoles.IsUnknown() || planRoles.IsNull() {
		if privateState != nil {
			diags.Append(privateState.SetKey(ctx, resourceSetImplicitRolesKey, []byte("[]"))...)
		}
		rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, backendRoles)
		diags.Append(setDiags...)
		return rolesSet, diags
	}

	var plannedRoles []string
	diags.Append(planRoles.ElementsAs(ctx, &plannedRoles, false)...)
	if diags.HasError() {
		return planRoles, diags
	}

	// Check if all planned roles are present in backendRoles.
	// If any planned role is missing, do not mask the discrepancy; clear implicit_roles and return backendRoles.
	for _, role := range plannedRoles {
		if !slices.Contains(backendRoles, role) {
			if privateState != nil {
				diags.Append(privateState.SetKey(ctx, resourceSetImplicitRolesKey, []byte("[]"))...)
			}
			rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, backendRoles)
			diags.Append(setDiags...)
			return rolesSet, diags
		}
	}

	// Compute implicitRoles = backendRoles \ plannedRoles
	var implicitRoles []string
	for _, role := range backendRoles {
		if !slices.Contains(plannedRoles, role) {
			implicitRoles = append(implicitRoles, role)
		}
	}

	if implicitRoles == nil {
		implicitRoles = []string{}
	}

	marshaled, err := json.Marshal(implicitRoles)
	if err != nil {
		diags.AddError("Error marshaling implicit roles", err.Error())
		return planRoles, diags
	}

	if privateState != nil {
		diags.Append(privateState.SetKey(ctx, resourceSetImplicitRolesKey, marshaled)...)
		if diags.HasError() {
			return planRoles, diags
		}
	}

	// Compute reconciledRoles = backendRoles \ implicitRoles
	var reconciledRoles []string
	for _, role := range backendRoles {
		if !slices.Contains(implicitRoles, role) {
			reconciledRoles = append(reconciledRoles, role)
		}
	}

	if reconciledRoles == nil {
		reconciledRoles = []string{}
	}

	rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, reconciledRoles)
	diags.Append(setDiags...)
	return rolesSet, diags
}

// filterRolesOnRead filters implicit system roles out of backend-returned roles during Read using private state data.
// If implicit_roles is absent, empty, or unmarshals to an empty slice, the full backendRoles slice is returned.
func filterRolesOnRead(
	ctx context.Context,
	backendRoles []string,
	privateData privateStateData,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	if privateData == nil {
		rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, backendRoles)
		diags.Append(setDiags...)
		return rolesSet, diags
	}

	rawBytes, getDiags := privateData.GetKey(ctx, resourceSetImplicitRolesKey)
	diags.Append(getDiags...)
	if diags.HasError() || len(rawBytes) == 0 {
		rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, backendRoles)
		diags.Append(setDiags...)
		return rolesSet, diags
	}

	var implicitRoles []string
	if err := json.Unmarshal(rawBytes, &implicitRoles); err != nil {
		diags.AddError("Error unmarshaling implicit roles", err.Error())
		return types.SetNull(types.StringType), diags
	}

	if len(implicitRoles) == 0 {
		rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, backendRoles)
		diags.Append(setDiags...)
		return rolesSet, diags
	}

	var filtered []string
	for _, role := range backendRoles {
		if !slices.Contains(implicitRoles, role) {
			filtered = append(filtered, role)
		}
	}

	if filtered == nil {
		filtered = []string{}
	}

	rolesSet, setDiags := types.SetValueFrom(ctx, types.StringType, filtered)
	diags.Append(setDiags...)
	return rolesSet, diags
}
