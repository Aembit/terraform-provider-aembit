package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// SignInPolicyResourceModel maps the resource schema.
type SignInPolicyResourceModel struct {
	// ID is required for Framework acceptance testing
	ID          types.String `tfsdk:"id"`
	SSORequired types.Bool   `tfsdk:"sso_required"`
	MFARequired types.Bool   `tfsdk:"mfa_required"`
}
