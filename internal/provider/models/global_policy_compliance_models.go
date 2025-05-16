package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type GlobalPolicyComplianceModel struct {
	// ID is required for Framework acceptance testing
	Id                             types.String `tfsdk:"id"`
	APTrustProviderCompliance      types.String `tfsdk:"access_policy_trust_provider_compliance"`
	APAccessConditionCompliance    types.String `tfsdk:"access_policy_access_condition_compliance"`
	ACTrustProviderCompliance      types.String `tfsdk:"agent_controller_trust_provider_compliance"`
	ACAllowedTLSHostanmeCompliance types.String `tfsdk:"agent_controller_allowed_tls_hostname_compliance"`
}
