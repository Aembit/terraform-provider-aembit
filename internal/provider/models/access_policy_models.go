package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// accessPolicyResourceModel maps the resource schema.
type AccessPolicyResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                  types.String                    `tfsdk:"id"`
	Name                types.String                    `tfsdk:"name"`
	IsActive            types.Bool                      `tfsdk:"is_active"`
	ClientWorkload      types.String                    `tfsdk:"client_workload"`
	TrustProviders      []types.String                  `tfsdk:"trust_providers"`
	AccessConditions    []types.String                  `tfsdk:"access_conditions"`
	CredentialProvider  types.String                    `tfsdk:"credential_provider"`
	CredentialProviders []*PolicyCredentialMappingModel `tfsdk:"credential_providers"`
	ServerWorkload      types.String                    `tfsdk:"server_workload"`
}

// accessPoliciesDataSourceModel maps the datasource schema.
type AccessPoliciesDataSourceModel struct {
	AccessPolicies []AccessPolicyResourceModel `tfsdk:"access_policies"`
}

type PolicyCredentialMappingModel struct {
	CredentialProviderId types.String `tfsdk:"credential_provider_id"`
	MappingType          types.String `tfsdk:"mapping_type"`
	AccountName          types.String `tfsdk:"account_name"`
	AccessKeyId 		 types.String `tfsdk:"access_key_id"`
	HeaderName           types.String `tfsdk:"header_name"`
	HeaderValue          types.String `tfsdk:"header_value"`
	HttpbodyFieldPath    types.String `tfsdk:"httpbody_field_path"`
	HttpbodyFieldValue   types.String `tfsdk:"httpbody_field_value"`
}
