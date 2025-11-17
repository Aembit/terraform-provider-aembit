package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.TrustProviderSecretResourceModel maps the resource schema.
type TrustProviderSecretResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String `tfsdk:"id"`
	TrustProviderID types.String `tfsdk:"trust_provider_id"`
	Secret          types.String `tfsdk:"secret"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Subject         types.String `tfsdk:"subject"`
	ExpiresAt       types.String `tfsdk:"expires_at"`
}
