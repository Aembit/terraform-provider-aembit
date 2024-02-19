package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// credentialProviderResourceModel maps the resource schema.
type credentialProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	IsActive               types.Bool   `tfsdk:"is_active"`
	ApiKey                 types.Object `tfsdk:"api_key"`
	OAuthClientCredentials types.Object `tfsdk:"oauth_client_credentials"`
	VaultClientToken       types.Object `tfsdk:"vault_client_token"`
}

// credentialProviderDataSourceModel maps the datasource schema.
type credentialProvidersDataSourceModel struct {
	CredentialProviders []credentialProviderResourceModel `tfsdk:"credential_providers"`
}

// credentialProviderApiKeyModel maps API Key credential configuration.
var credentialProviderApiKeyModel = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"api_key": types.StringType,
	},
}

// credentialProviderApiKeyModel maps OAuth Client Credentials Flow configuration.
var credentialProviderOAuthClientCredentialsModel = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"token_url":     types.StringType,
		"client_id":     types.StringType,
		"client_secret": types.StringType,
		"scopes":        types.StringType,
	},
}

// credentialProviderApiKeyModel maps OAuth Client Credentials Flow configuration.
var credentialProviderVaultClientTokenModel = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"subject":      types.StringType,
		"subject_type": types.StringType,
		"custom_claims": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":        types.StringType,
					"value":      types.StringType,
					"value_type": types.StringType,
				},
			},
		},
		"lifetime":         types.Int64Type,
		"vault_host":       types.StringType,
		"vault_tls":        types.BoolType,
		"vault_port":       types.Int64Type,
		"vault_namespace":  types.StringType,
		"vault_role":       types.StringType,
		"vault_path":       types.StringType,
		"vault_forwarding": types.StringType,
	},
}
