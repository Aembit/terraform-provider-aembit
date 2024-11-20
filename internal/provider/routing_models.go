package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// routingListResourceModel maps the resource schema.
type routingListResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	ProxyUrl        types.String `tfsdk:"proxy_url"`
	ResourceSetName types.String `tfsdk:"resource_set_name"`
}

type routingResourceModel struct {
	// ID is required for Framework acceptance testing
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	ProxyUrl      types.String `tfsdk:"proxy_url"`
	ResourceSetId types.String `tfsdk:"resource_set_id"`
}

// routingDataSourceModel maps the datasource schema.
type routingsDataSourceModel struct {
	Routings []routingListResourceModel `tfsdk:"routings"`
}
