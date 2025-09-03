package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentityProviderResourceModel struct {
	ID                        types.String                `tfsdk:"id"`
	Name                      types.String                `tfsdk:"name"`
	Description               types.String                `tfsdk:"description"`
	IsActive                  types.Bool                  `tfsdk:"is_active"`
	Tags                      types.Map                   `tfsdk:"tags"`
	MetadataUrl               types.String                `tfsdk:"metadata_url"`
	MetadataXml               types.String                `tfsdk:"metadata_xml"`
	SamlStatementRoleMappings []SamlStatementRoleMappings `tfsdk:"saml_statement_role_mappings"`
}

type SamlStatementRoleMappings struct {
	AttributeName  types.String   `tfsdk:"attribute_name"`
	AttributeValue types.String   `tfsdk:"attribute_value"`
	Roles          []types.String `tfsdk:"roles"`
}

type IdentityProviderDataSourceModel struct {
	IdentityProviders []IdentityProviderResourceModel `tfsdk:"identity_providers"`
}
