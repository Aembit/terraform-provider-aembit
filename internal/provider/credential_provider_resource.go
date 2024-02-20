package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	state = convertCredentialProviderDTOToModel(ctx, credential_provider, state)

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

func convertCredentialProviderDTOToModel(ctx context.Context, dto aembit.CredentialProviderDTO, state credentialProviderResourceModel) credentialProviderResourceModel {
	var model credentialProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalId)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)

	// Set the objects to null to begin with
	model.ApiKey = types.ObjectNull(credentialProviderApiKeyModel.AttrTypes)
	model.OAuthClientCredentials = types.ObjectNull(credentialProviderOAuthClientCredentialsModel.AttrTypes)
	model.VaultClientToken = types.ObjectNull(credentialProviderVaultClientTokenModel.AttrTypes)

	// Now fill in the objects based on the Credential Provider type
	switch dto.Type {
	case "apikey":
		// We leave the API Key alone complete to avoid confusing Terraform since Aembit doesn't return the current secret
		model.ApiKey = convertApiKeyDTOToModel(ctx, dto, state)
	case "oauth-client-credential":
		// Generally, we want to pass the same state object through so that no change is detected, unless we see
		// that a non-secret value has changed. So let's compare values, and only update if there has been any change.

		// First, parse the credential_provider.ProviderDetail JSON returned from Aembit Cloud
		var oauth aembit.CredentialOAuthClientCredentialDTO
		json.Unmarshal([]byte(dto.ProviderDetail), &oauth)

		// Compare for changes
		if getStringAttr(ctx, state.OAuthClientCredentials, "token_url") != oauth.TokenUrl ||
			getStringAttr(ctx, state.OAuthClientCredentials, "client_id") != oauth.ClientID ||
			getStringAttr(ctx, state.OAuthClientCredentials, "scopes") != oauth.Scope {

			// Pull the client_secret from the state to avoid confusing Terraform since Aembit doesn't return the current secret
			var secret string = getStringAttr(ctx, state.OAuthClientCredentials, "client_secret")

			model.OAuthClientCredentials, _ = types.ObjectValue(credentialProviderOAuthClientCredentialsModel.AttrTypes,
				map[string]attr.Value{
					"token_url":     types.StringValue(oauth.TokenUrl),
					"client_id":     types.StringValue(oauth.ClientID),
					"client_secret": types.StringValue(secret),
					"scopes":        types.StringValue(oauth.Scope),
				})
		} else {
			model.OAuthClientCredentials = state.OAuthClientCredentials
		}
	case "vaultClientToken":
		// First, parse the credential_provider.ProviderDetail JSON returned from Aembit Cloud
		var vault aembit.CredentialVaultClientTokenDTO
		json.Unmarshal([]byte(dto.ProviderDetail), &vault)

		// Get the custom claims to be injected into the model
		claims := make([]attr.Value, len(vault.JwtConfig.CustomClaims))
		//types.ObjectValue(credentialProviderVaultClientTokenCustomClaimsModel.AttrTypes),
		//claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
		for i, claim := range vault.JwtConfig.CustomClaims {
			claims[i], _ = types.ObjectValue(credentialProviderVaultClientTokenCustomClaimsModel.AttrTypes,
				map[string]attr.Value{
					"key":        types.StringValue(claim.Key),
					"value":      types.StringValue(claim.Value),
					"value_type": types.StringValue(claim.ValueType),
				})
		}
		claimsValue, _ := types.SetValue(credentialProviderVaultClientTokenCustomClaimsModel, claims)

		// Construct the model
		model.VaultClientToken, _ = types.ObjectValue(credentialProviderVaultClientTokenModel.AttrTypes,
			map[string]attr.Value{
				"subject":          types.StringValue(vault.JwtConfig.Subject),
				"subject_type":     types.StringValue(vault.JwtConfig.SubjectType),
				"custom_claims":    claimsValue,
				"lifetime":         types.Int64Value(int64(vault.JwtConfig.Lifetime)),
				"vault_host":       types.StringValue(vault.VaultCluster.VaultHost),
				"vault_tls":        types.BoolValue(vault.VaultCluster.Tls),
				"vault_port":       types.Int64Value(int64(vault.VaultCluster.Port)),
				"vault_namespace":  types.StringValue(vault.VaultCluster.Namespace),
				"vault_role":       types.StringValue(vault.VaultCluster.Role),
				"vault_path":       types.StringValue(vault.VaultCluster.AuthenticationPath),
				"vault_forwarding": types.StringValue(vault.VaultCluster.ForwardingConfig),
			})
	}
	return model
}

func convertApiKeyDTOToModel(ctx context.Context, dto aembit.CredentialProviderDTO, state credentialProviderResourceModel) types.ObjectValue {
	value, _ := types.ObjectValue(credentialProviderApiKeyModel.AttrTypes,
		map[string]attr.Value{
			"api_key": types.StringValue(getStringAttr(ctx, state.ApiKey, "api_key")),
		})
	return value
}
