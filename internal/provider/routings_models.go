package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// routingResourceModel maps the resource schema.
type routingResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	IsActive        types.Bool   `tfsdk:"is_active"`
	ProxyUrl        types.String `tfsdk:"proxy_url"`
	ResourceSetId   types.String `tfsdk:"resource_set_id"`
	ResourceSetName types.String `tfsdk:"resource_set_id"`
}

// routingDataSourceModel maps the datasource schema.
type routingsDataSourceModel struct {
	Routings []routingResourceModel `tfsdk:"routings"`
}
