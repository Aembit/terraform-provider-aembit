package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CredentialProviderIntegrationResourceModel maps the resource schema.
type CredentialProviderIntegrationResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                   types.String                                            `tfsdk:"id"`
	Name                 types.String                                            `tfsdk:"name"`
	Description          types.String                                            `tfsdk:"description"`
	GitLab               *CredentialProviderIntegrationGitlabModel               `tfsdk:"gitlab"`
	AwsIamRole           *CredentialProviderIntegrationAwsIamRoleModel           `tfsdk:"aws_iam_role"`
	AzureEntraFederation *CredentialProviderIntegrationAzureEntraFederationModel `tfsdk:"azure_entra_federation"`
}

// CredentialProviderIntegrationsDataSourceModel maps the datasource schema.
type CredentialProviderIntegrationsDataSourceModel struct {
	CredentialProviderIntegrations []CredentialProviderIntegrationResourceModel `tfsdk:"credential_provider_integrations"`
}

type CredentialProviderIntegrationGitlabModel struct {
	Url                 types.String `tfsdk:"url"`
	PersonalAccessToken types.String `tfsdk:"personal_access_token"`
}

type CredentialProviderIntegrationAwsIamRoleModel struct {
	RoleArn           types.String `tfsdk:"role_arn"`
	LifetimeInSeconds types.Int32  `tfsdk:"lifetime_in_seconds"`
	OIDCIssuerUrl     types.String `tfsdk:"oidc_issuer_url"`
	TokenAudience     types.String `tfsdk:"token_audience"`
}

type CredentialProviderIntegrationAzureEntraFederationModel struct {
	Audience      types.String `tfsdk:"audience"`
	Subject       types.String `tfsdk:"subject"`
	AzureTenant   types.String `tfsdk:"azure_tenant"`
	ClientID      types.String `tfsdk:"client_id"`
	KeyVaultName  types.String `tfsdk:"key_vault_name"`
	OIDCIssuerUrl types.String `tfsdk:"oidc_issuer_url"`
}
