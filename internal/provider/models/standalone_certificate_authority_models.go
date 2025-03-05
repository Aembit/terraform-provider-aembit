package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// StandaloneCertificateResourceModel maps the resource schema.
type StandaloneCertificateAuthorityResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Tags                types.Map    `tfsdk:"tags"`
	LeafLifetime        types.Int32  `tfsdk:"leaf_lifetime"`
	NotBefore           types.String `tfsdk:"not_before"`
	NotAfter            types.String `tfsdk:"not_after"`
	ClientWorkloadCount types.Int32  `tfsdk:"client_workload_count"`
	ResourceSetCount    types.Int32  `tfsdk:"resource_set_count"`
}

// StandaloneCertificatesDataSourceModel maps the datasource schema.
type StandaloneCertificateAuthoritiesDataSourceModel struct {
	StandaloneCertificates []StandaloneCertificateAuthorityResourceModel `tfsdk:"standalone_certificate_authorities"`
}
