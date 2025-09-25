package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// aembitProviderModel maps provider schema data to a Go type.
type AembitProviderrModel struct {
	Tenant      types.String `tfsdk:"tenant"`
	Token       types.String `tfsdk:"token"`
	ClientID    types.String `tfsdk:"client_id"`
	ResourceSet types.String `tfsdk:"resource_set_id"`
	DefaultTags types.Object `tfsdk:"default_tags"`
}
