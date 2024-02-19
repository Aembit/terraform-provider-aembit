package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
			"vault_client_token": schema.ObjectAttribute{
				Optional:       true,
				AttributeTypes: credentialProviderVaultClientTokenModel.AttrTypes,
			},
		},
	}
}

// Configure validators to ensure that only one credential provider type is specified
func (r *credentialProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("api_key"),
			path.MatchRoot("oauth_client_credentials"),
			path.MatchRoot("vault_client_token"),
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
	var credential aembit.CredentialProviderDTO = convertCredentialProviderModelToDTO(ctx, plan, nil)

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
	case "apikey":
	// We leave the API Key alone complete to avoid confusing Terraform since Aembit doesn't return the current secret
	case "oauth-client-credential":
		// Generally, we want to pass the same state object through so that no change is detected, unless we see
		// that a non-secret value has changed. So let's compare values, and only update if there has been any change.

		// First, parse the credential_provider.ProviderDetail JSON returned from Aembit Cloud
		var oauth aembit.CredentialOAuthClientCredentialDTO
		json.Unmarshal([]byte(credential_provider.ProviderDetail), &oauth)

		// Compare for changes
		if strings.Trim(state.OAuthClientCredentials.Attributes()["token_url"].String(), "\"") != oauth.TokenUrl ||
			strings.Trim(state.OAuthClientCredentials.Attributes()["client_id"].String(), "\"") != oauth.ClientID ||
			strings.Trim(state.OAuthClientCredentials.Attributes()["scopes"].String(), "\"") != oauth.Scope {
			fmt.Printf("TokenURL: %s %s\n", strings.Trim(state.OAuthClientCredentials.Attributes()["token_url"].String(), "\""), oauth.TokenUrl)
			fmt.Printf("TokenURL: %s %s\n", strings.Trim(state.OAuthClientCredentials.Attributes()["client_id"].String(), "\""), oauth.ClientID)
			fmt.Printf("TokenURL: %s %s\n", strings.Trim(state.OAuthClientCredentials.Attributes()["scopes"].String(), "\""), oauth.Scope)

			// Pull the client_secret from the state to avoid confusing Terraform since Aembit doesn't return the current secret
			var secret string = state.OAuthClientCredentials.Attributes()["client_secret"].String()

			state.OAuthClientCredentials, _ = types.ObjectValue(credentialProviderOAuthClientCredentialsModel.AttrTypes,
				map[string]attr.Value{
					"token_url":     types.StringValue(oauth.TokenUrl),
					"client_id":     types.StringValue(oauth.ClientID),
					"client_secret": types.StringValue(secret),
					"scopes":        types.StringValue(oauth.Scope),
				})
		}
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
	var credential aembit.CredentialProviderDTO = convertCredentialProviderModelToDTO(ctx, plan, &external_id)

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

func convertCredentialProviderModelToDTO(ctx context.Context, model credentialProviderResourceModel, external_id *string) aembit.CredentialProviderDTO {
	var credential aembit.CredentialProviderDTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if external_id != nil {
		credential.EntityDTO.ExternalId = *external_id
	}

	// Handle the API Key use case
	if !model.ApiKey.IsNull() {
		credential.Type = "apikey"
		apiKey := aembit.CredentialApiKeyDTO{ApiKey: getStringAttr(ctx, model.ApiKey, "api_key")}
		apiKeyJson, _ := json.Marshal(apiKey)
		credential.ProviderDetail = string(apiKeyJson)
	}

	// Handle the OAuth Client Credentials use case
	if !model.OAuthClientCredentials.IsNull() {
		credential.Type = "oauth-client-credential"
		oauth := aembit.CredentialOAuthClientCredentialDTO{
			TokenUrl:        getStringAttr(ctx, model.OAuthClientCredentials, "token_url"),
			ClientID:        getStringAttr(ctx, model.OAuthClientCredentials, "client_id"),
			ClientSecret:    getStringAttr(ctx, model.OAuthClientCredentials, "client_secret"),
			Scope:           getStringAttr(ctx, model.OAuthClientCredentials, "scopes"),
			CredentialStyle: "authHeader",
		}
		oauthJson, _ := json.Marshal(oauth)
		credential.ProviderDetail = string(oauthJson)
	}

	// Handle the Vault Cvlient Token use case
	if !model.VaultClientToken.IsNull() {
		credential.Type = "vaultClientToken"
		vault := aembit.CredentialVaultClientTokenDTO{
			JwtConfig: &aembit.CredentialVaultClientTokenJwtConfigDTO{
				Issuer:       "https://62c41c.id.aembit.local/",
				Subject:      getStringAttr(ctx, model.VaultClientToken, "subject"),
				SubjectType:  getStringAttr(ctx, model.VaultClientToken, "subject_type"),
				Lifetime:     getInt32Attr(ctx, model.VaultClientToken, "lifetime"),
				CustomClaims: make([]aembit.CredentialVaultClientTokenClaimsDTO, 0),
			},
			VaultCluster: &aembit.CredentialVaultClientTokenVaultClusterDTO{
				VaultHost:          getStringAttr(ctx, model.VaultClientToken, "vault_host"),
				Port:               getInt32Attr(ctx, model.VaultClientToken, "vault_port"),
				Tls:                getBoolAttr(ctx, model.VaultClientToken, "vault_tls"),
				Namespace:          getStringAttr(ctx, model.VaultClientToken, "vault_namespace"),
				Role:               getStringAttr(ctx, model.VaultClientToken, "vault_role"),
				AuthenticationPath: getStringAttr(ctx, model.VaultClientToken, "vault_path"),
				ForwardingConfig:   getStringAttr(ctx, model.VaultClientToken, "vault_forwarding"),
			},
		}
		claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
		for _, claim := range claims {
			vault.JwtConfig.CustomClaims = append(vault.JwtConfig.CustomClaims, aembit.CredentialVaultClientTokenClaimsDTO{
				Key:       getStringAttr(ctx, claim, "key"),
				Value:     getStringAttr(ctx, claim, "value"),
				ValueType: getStringAttr(ctx, claim, "value_type"),
			})
		}
		vaultJson, _ := json.Marshal(vault)
		credential.ProviderDetail = string(vaultJson)
	}
	return credential
}

func getStringAttr(ctx context.Context, tfObject types.Object, name string) string {
	var value string
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	return value
}

func getInt32Attr(ctx context.Context, tfObject types.Object, name string) int32 {
	var value big.Float
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	result, _ := value.Int64()
	return int32(result)
}

func getBoolAttr(ctx context.Context, tfObject types.Object, name string) bool {
	var value bool
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	return value
}

func getSetObjectAttr(ctx context.Context, tfObject types.Object, name string) []types.Object {
	var objSlice []types.Object
	tfsdk.ValueAs(ctx, tfObject.Attributes()[name], &objSlice)

	objects := make([]types.Object, len(objSlice))
	for i, val := range objSlice {
		objects[i] = val
	}

	return objects
}
