package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// trustProviderResourceModel maps the resource schema.
type trustProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID            types.String                     `tfsdk:"id"`
	Name          types.String                     `tfsdk:"name"`
	Description   types.String                     `tfsdk:"description"`
	IsActive      types.Bool                       `tfsdk:"is_active"`
	AzureMetadata *trustProviderAzureMetadataModel `tfsdk:"azure_metadata"`
}

// trustProviderDataSourceModel maps the datasource schema.
type trustProvidersDataSourceModel struct {
	TrustProviders []trustProviderResourceModel `tfsdk:"trust_providers"`
}

type trustProviderAzureMetadataModel struct {
	Sku            types.String `tfsdk:"sku"`
	VmId           types.String `tfsdk:"vm_id"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
}
