package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// credentialProviderResourceModel maps the resource schema.
type credentialProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsActive    types.Bool   `tfsdk:"is_active"`
	ApiKey      types.Object `tfsdk:"api_key"`
	Type        types.String `tfsdk:"type"`
}

// credentialProviderDataSourceModel maps the datasource schema.
type credentialProvidersDataSourceModel struct {
	CredentialProviders []credentialProviderResourceModel `tfsdk:"credential_providers"`
}

// serviceEndpointModel maps service endpoint data.
type credentialProviderApiKeyModel struct {
	Value types.String `tfsdk:"value"`
}
