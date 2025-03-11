package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.ResourceSetResourceModel maps the resource schema.
type ResourceSetResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                             types.String   `tfsdk:"id"`
	Name                           types.String   `tfsdk:"name"`
	Description                    types.String   `tfsdk:"description"`
	Roles                          []types.String `tfsdk:"roles"`
	StandaloneCertificateAuthority types.String   `tfsdk:"standalone_certificate_authority"`
}

// models.ResourceSetsDataModel maps the list of all resource sets.
type ResourceSetsDataModel struct {
	ResourceSets []ResourceSetResourceModel `tfsdk:"resource_sets"`
}
