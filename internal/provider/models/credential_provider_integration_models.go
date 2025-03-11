package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialProviderIntegrationResourceModel maps the resource schema.
type CredentialProviderIntegrationResourceModel struct {
	// ID is required for Framework acceptance testing
	ID          types.String                              `tfsdk:"id"`
	Name        types.String                              `tfsdk:"name"`
	Description types.String                              `tfsdk:"description"`
	GitLab      *CredentialProviderIntegrationGitlabModel `tfsdk:"gitlab"`
}

// CredentialProviderIntegrationsDataSourceModel maps the datasource schema.
type CredentialProviderIntegrationsDataSourceModel struct {
	CredentialProviderIntegrations []CredentialProviderIntegrationResourceModel `tfsdk:"credential_provider_integrations"`
}

type CredentialProviderIntegrationGitlabModel struct {
	Url                 types.String `tfsdk:"url"`
	PersonalAccessToken types.String `tfsdk:"personal_access_token"`
}
