package provider

import (
	"context"
	"fmt"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &standaloneCertificateAuthorityResource{}
	_ resource.ResourceWithConfigure   = &standaloneCertificateAuthorityResource{}
	_ resource.ResourceWithImportState = &standaloneCertificateAuthorityResource{}
	_ resource.ResourceWithModifyPlan  = &standaloneCertificateAuthorityResource{}
)

// NewServerWorkloadResource is a helper function to simplify the provider implementation.
func NewStandaloneCertificateAuthorityResource() resource.Resource {
	return &standaloneCertificateAuthorityResource{}
}

// serverWorkloadResource is the resource implementation.
type standaloneCertificateAuthorityResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *standaloneCertificateAuthorityResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_standalone_certificate_authority"
}

// Configure adds the provider configured client to the resource.
func (r *standaloneCertificateAuthorityResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *standaloneCertificateAuthorityResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the standalone certificate authority.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "User-provided name for the standalone certificate authority.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "User-provided description for the standalone certificate authority.",
				Optional:    true,
				Computed:    true,
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"leaf_lifetime": schema.Int32Attribute{
				Description: "Leaf certificate lifetime(in minutes) of the standalone certificate authority. Valid options; 60, 1440, 10080.",
				Required:    true,
				Validators: []validator.Int32{
					int32validator.OneOf(60, 1440, 10080),
				},
			},
			"not_before": schema.StringAttribute{
				Description: "ISO 8601 formatted not before date of the standalone certificate authority.",
				Computed:    true,
			},
			"not_after": schema.StringAttribute{
				Description: "ISO 8601 formatted not after date of the standalone certificate authority.",
				Computed:    true,
			},
			"client_workload_count": schema.Int32Attribute{
				Description: "Client workloads associated with the standalone certificate authority.",
				Computed:    true,
			},
			"resource_set_count": schema.Int32Attribute{
				Description: "Resource sets associated with the standalone certificate authority.",
				Computed:    true,
			},
		},
	}
}

func (r *standaloneCertificateAuthorityResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
}

// Create creates the resource and sets the initial Terraform state.
func (r *standaloneCertificateAuthorityResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.StandaloneCertificateAuthorityResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	standaloneCertificate := convertStandaloneCertificateModelToDTO(
		ctx,
		plan,
		nil,
		r.client.DefaultTags,
	)

	// Create new Standalone Certificate Authority
	standaloneCertificateResponse, err := r.client.CreateStandaloneCertificate(
		standaloneCertificate,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Standalone Certificate Authority",
			"Could not create Standalone Certificate Authority, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertStandaloneCertificateDTOToModel(ctx, *standaloneCertificateResponse, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *standaloneCertificateAuthorityResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.StandaloneCertificateAuthorityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed workload value from Aembit
	standaloneCertificate, err, notFound := r.client.GetStandaloneCertificate(
		state.ID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Standalone Certificate Authority",
			fmt.Sprintf(
				"An error occurred while attempting to fetch the Standalone Certificate Authority with External ID '%s' from Terraform state: %v",
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
	state = convertStandaloneCertificateDTOToModel(ctx, standaloneCertificate, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *standaloneCertificateAuthorityResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.StandaloneCertificateAuthorityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.StandaloneCertificateAuthorityResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	workload := convertStandaloneCertificateModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)

	// Update Standalone Certificate Authority
	serverWorkload, err := r.client.UpdateStandaloneCertificate(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Standalone Certificate Authority",
			"Could not update Standalone Certificate Authority, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertStandaloneCertificateDTOToModel(ctx, *serverWorkload, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *standaloneCertificateAuthorityResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.StandaloneCertificateAuthorityResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Standalone Certificate Authority
	_, err := r.client.DeleteStandaloneCertificate(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Standalone Certificate Authority",
			"Could not delete Standalone Certificate Authority, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalID.
func (r *standaloneCertificateAuthorityResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertStandaloneCertificateModelToDTO(
	ctx context.Context,
	model models.StandaloneCertificateAuthorityResourceModel,
	externalID *string,
	defaultTags map[string]string,
) aembit.StandaloneCertificateDTO {
	var standaloneCertificate aembit.StandaloneCertificateDTO
	standaloneCertificate.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
	}

	standaloneCertificate.LeafLifetime = model.LeafLifetime.ValueInt32()
	standaloneCertificate.NotBefore = model.NotBefore.ValueString()
	standaloneCertificate.NotAfter = model.NotAfter.ValueString()

	if externalID != nil {
		standaloneCertificate.ExternalID = *externalID
	}

	standaloneCertificate.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)
	return standaloneCertificate
}

func convertStandaloneCertificateDTOToModel(
	ctx context.Context,
	dto aembit.StandaloneCertificateDTO,
	planModel *models.StandaloneCertificateAuthorityResourceModel,
) models.StandaloneCertificateAuthorityResourceModel {
	var model models.StandaloneCertificateAuthorityResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.LeafLifetime = types.Int32Value(dto.LeafLifetime)
	model.NotBefore = types.StringValue(dto.NotBefore)
	model.NotAfter = types.StringValue(dto.NotAfter)
	model.ClientWorkloadCount = types.Int32Value(dto.ClientWorkloadCount)
	model.ResourceSetCount = types.Int32Value(dto.ResourceSetCount)
	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	return model
}
