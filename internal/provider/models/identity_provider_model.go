package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentityProviderResourceModel struct {
	ID                       types.String               `tfsdk:"id"`
	Name                     types.String               `tfsdk:"name"`
	Description              types.String               `tfsdk:"description"`
	IsActive                 types.Bool                 `tfsdk:"is_active"`
	Tags                     types.Map                  `tfsdk:"tags"`
	TagsAll                  types.Map                  `tfsdk:"tags_all"`
	SsoStatementRoleMappings []SsoStatementRoleMappings `tfsdk:"sso_statement_role_mappings"`
	Saml                     *IdentityProviderSamlModel `tfsdk:"saml"`
	Oidc                     *IdentityProviderOidcModel `tfsdk:"oidc"`
}

type IdentityProviderSamlModel struct {
	MetadataUrl             types.String `tfsdk:"metadata_url"`
	MetadataXml             types.String `tfsdk:"metadata_xml"`
	ServiceProviderEntityId types.String `tfsdk:"service_provider_entity_id"`
	ServiceProviderSsoUrl   types.String `tfsdk:"service_provider_sso_url"`
}

type IdentityProviderOidcModel struct {
	OidcBaseUrl       types.String `tfsdk:"oidc_base_url"`
	ClientId          types.String `tfsdk:"client_id"`
	Scopes            types.String `tfsdk:"scopes"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	AuthType          types.String `tfsdk:"auth_type"`
	PkceRequired      types.Bool   `tfsdk:"pcke_required"`
	AembitRedirectUrl types.String `tfsdk:"aembit_redirect_url"`
	AembitJwksUrl     types.String `tfsdk:"aembit_jwks_url"`
}

type SsoStatementRoleMappings struct {
	AttributeName  types.String   `tfsdk:"attribute_name"`
	AttributeValue types.String   `tfsdk:"attribute_value"`
	Roles          []types.String `tfsdk:"roles"`
}

type IdentityProviderDataSourceModel struct {
	IdentityProviders []IdentityProviderResourceModel `tfsdk:"identity_providers"`
}
