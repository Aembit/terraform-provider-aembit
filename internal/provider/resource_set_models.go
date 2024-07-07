package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourcesetResourceModel maps the resource schema.
type resourceSetResourceModel struct {
	// ID is required for Framework acceptance testing
	ID          types.String   `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Roles       []types.String `tfsdk:"roles"`
}

// resourceSetsDataModel maps the list of all resource sets.
type resourceSetsDataModel struct {
	ResourceSets []resourceSetResourceModel `tfsdk:"resource_sets"`
}
