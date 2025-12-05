package provider

import (
	"context"
	"encoding/base64"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &trustProviderSecretResource{}
	_ resource.ResourceWithConfigure   = &trustProviderSecretResource{}
	_ resource.ResourceWithImportState = &trustProviderSecretResource{}
)

// NewTrustProviderSecret is a helper function to simplify the provider implementation.
func NewTrustProviderSecretResource() resource.Resource {
	return &trustProviderSecretResource{}
}

// trustProviderSecretResource is the resource implementation.
type trustProviderSecretResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *trustProviderSecretResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_trust_provider_secret"
}

// Configure adds the provider configured client to the resource.
func (r *trustProviderSecretResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *trustProviderSecretResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Trust Provider Secret.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"trust_provider_id": schema.StringAttribute{
				Description: "Unique identifier of the Trust Provider.",
				Required:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"secret": schema.StringAttribute{
				Description: "PEM Certificate or Symmetric Key to be used for Signature verification.",
				Required:    true,
				Sensitive:   true,
			},
			"name": schema.StringAttribute{
				Description: "Thumbprint of the secret.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the Secret. Possible values are: \n" +
					"\t* `Certificate`\n" +
					"\t* `SymmetricKey`\n",
				Optional: true,
				Computed: true,
				Default: stringdefault.StaticString(
					"Certificate",
				),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"Certificate",
						"SymmetricKey",
					}...),
				},
			},
			"subject": schema.StringAttribute{
				Description: "Subject of the Certificate.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "Expiration date of the certificate.",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *trustProviderSecretResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.TrustProviderSecretResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	tpSecret := convertTrustProviderSecretResourceModelToDTO(plan, nil)

	// Create new Trust Provider Secret
	trustProviderSecret, err := r.client.UpsertTrustProviderSecret(
		tpSecret,
		plan.TrustProviderID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Provider Secret",
			"Could not create Trust Provider Secret, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderSecretDTOToModel(*trustProviderSecret, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trustProviderSecretResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.TrustProviderSecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed tpSecret value from Aembit
	trustProviderSecret, err, notFound := r.client.GetTrustProviderSecret(
		state.TrustProviderID.ValueString(),
		state.ID.ValueString(),
		nil,
	)

	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Trust Provider Secret",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertTrustProviderSecretDTOToModel(trustProviderSecret, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustProviderSecretResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.TrustProviderSecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.TrustProviderSecretResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	tpSecret := convertTrustProviderSecretResourceModelToDTO(plan, &externalID)

	// Update Trust Provider Secret
	trustProviderSecret, err := r.client.UpsertTrustProviderSecret(
		tpSecret,
		plan.TrustProviderID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trust Provider Secret",
			"Could not update Trust Provider Secret, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertTrustProviderSecretDTOToModel(*trustProviderSecret, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *trustProviderSecretResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.TrustProviderSecretResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Trust Provider Secret
	_, err := r.client.DeleteTrustProviderSecret(
		ctx,
		state.TrustProviderID.ValueString(),
		state.ID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider Secret",
			"Could not delete Trust Provider Secret, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *trustProviderSecretResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertTrustProviderSecretResourceModelToDTO(
	model models.TrustProviderSecretResourceModel,
	externalID *string,
) aembit.TrustProviderSecretDTO {
	var tpSecret aembit.TrustProviderSecretDTO
	tpSecret.Name = model.Name.ValueString()
	tpSecret.Subject = model.Subject.ValueString()
	tpSecret.Type = model.Type.ValueString()
	tpSecret.ExpiresAt = model.ExpiresAt.ValueString()
	tpSecret.Secret = base64.StdEncoding.EncodeToString(
		[]byte(model.Secret.ValueString()),
	)

	if externalID != nil {
		tpSecret.ExternalID = *externalID
	}
	return tpSecret
}

func convertTrustProviderSecretDTOToModel(
	dto aembit.TrustProviderSecretDTO,
	planModel *models.TrustProviderSecretResourceModel,
) models.TrustProviderSecretResourceModel {
	var model models.TrustProviderSecretResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Subject = types.StringValue(dto.Subject)
	model.Type = types.StringValue(dto.Type)
	model.ExpiresAt = types.StringValue(dto.ExpiresAt)

	model.Secret = types.StringNull()
	if planModel != nil {
		model.Secret = planModel.Secret
		model.TrustProviderID = planModel.TrustProviderID
	}

	return model
}
