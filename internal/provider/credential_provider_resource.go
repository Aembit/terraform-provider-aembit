package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &credentialProviderResource{}
	_ resource.ResourceWithConfigure   = &credentialProviderResource{}
	_ resource.ResourceWithImportState = &credentialProviderResource{}
)

// NewCredentialProviderResource is a helper function to simplify the provider implementation.
func NewCredentialProviderResource() resource.Resource {
	return &credentialProviderResource{}
}

// credentialProviderResource is the resource implementation.
type credentialProviderResource struct {
	client *aembit.Client
}

// Metadata returns the resource type name.
func (r *credentialProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider"
}

// Configure adds the provider configured client to the resource.
func (r *credentialProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *credentialProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Alphanumeric identifier of the credential provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "User-provided name of the credential provider.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "User-provided description of the credential provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the credential provider.",
				Optional:    true,
				Computed:    true,
			},
			"api_key": schema.ObjectAttribute{
				Optional:       true,
				Sensitive:      true,
				AttributeTypes: credentialProviderApiKeyModel.AttrTypes,
			},
			"oauth_client_credentials": schema.ObjectAttribute{
				Optional:       true,
				Sensitive:      true,
				AttributeTypes: credentialProviderOAuthClientCredentialsModel.AttrTypes,
			},
		},
	}
}

func (r *credentialProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("api_key"),
			path.MatchRoot("oauth_client_credentials"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *credentialProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderDTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}
	if !plan.ApiKey.IsNull() {
		credential.Type = "apikey"
		credential.ProviderDetail = "{\"ApiKey\":\"test\"}"
	} else if !plan.OAuthClientCredentials.IsNull() {
		credential.Type = "oauth-client-credential"
		credential.ProviderDetail = "{ \"Url\": \"https://aembit.io/token\", \"ClientID\": \"testr\", \"ClientSecret\": \"testrt\", \"Scope\": \"test\", \"CredentialStyle\": \"authHeader\", \"CustomParameters\" : [] }"
	}

	// Create new Credential Provider
	credential_provider, err := r.client.CreateCredentialProvider(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating credential provider",
			"Could not create credential provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(credential_provider.EntityDTO.ExternalId)
	plan.Name = types.StringValue(credential_provider.EntityDTO.Name)
	plan.Description = types.StringValue(credential_provider.EntityDTO.Description)
	plan.IsActive = types.BoolValue(credential_provider.EntityDTO.IsActive)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *credentialProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed credential value from Aembit
	credential_provider, err := r.client.GetCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Aembit Credential Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(credential_provider.EntityDTO.ExternalId)
	state.Name = types.StringValue(credential_provider.EntityDTO.Name)
	state.Description = types.StringValue(credential_provider.EntityDTO.Description)
	state.IsActive = types.BoolValue(credential_provider.EntityDTO.IsActive)

	// Make sure the type and non-secret values have not changed
	switch credential_provider.Type {
	//case "oauth-client-credential":
	//	state.OAuthClientCredentials, _ = types.ObjectValue(credentialProviderOAuthClientCredentialsModel.AttrTypes,
	//		map[string]attr.Value{
	//			"token_url":     types.StringValue("test"),
	//			"client_id":     types.StringValue("test"),
	//			"client_secret": types.StringValue(""),
	//			"scopes":        types.StringValue("test"),
	//		})
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *credentialProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	var external_id string
	external_id = state.ID.ValueString()

	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderDTO
	credential.EntityDTO = aembit.EntityDTO{
		ExternalId:  external_id,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsActive:    plan.IsActive.ValueBool(),
	}
	credential.Type = "apikey"
	credential.ProviderDetail = "{\"ApiKey\":\"test\"}"

	// Update Credential Provider
	credential_provider, err := r.client.UpdateCredentialProvider(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating credential provider",
			"Could not update credential provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(credential_provider.EntityDTO.ExternalId)
	plan.Name = types.StringValue(credential_provider.EntityDTO.Name)
	plan.Description = types.StringValue(credential_provider.EntityDTO.Description)
	plan.IsActive = types.BoolValue(credential_provider.EntityDTO.IsActive)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *credentialProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Credential Provider is Active
	if state.IsActive == types.BoolValue(true) {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Credential Provider is active and cannot be deleted. Please mark the credential as inactive first.",
		)
		return
	}

	// Delete existing Credential Provider
	_, err := r.client.DeleteCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Could not delete credential provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId
func (r *credentialProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
