package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.ClientWorkloadResourceModel maps the resource schema.
type ClientWorkloadResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	IsActive                       types.Bool   `tfsdk:"is_active"`
	Identities                     types.Set    `tfsdk:"identities"`
	Tags                           types.Map    `tfsdk:"tags"`
	TagsAll                        types.Map    `tfsdk:"tags_all"`
	StandaloneCertificateAuthority types.String `tfsdk:"standalone_certificate_authority"`
}

// clientWorkloadDataSourceModel maps the datasource schema.
type ClientWorkloadsDataSourceModel struct {
	ClientWorkloads []ClientWorkloadResourceModel `tfsdk:"client_workloads"`
}

// models.IdentitiesModel maps client workload identity data.
type IdentitiesModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// TfIdentityObjectType maps client workload identity data to an Object type.
var TfIdentityObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type":  types.StringType,
	"value": types.StringType,
}}
