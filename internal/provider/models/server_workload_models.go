package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ServerWorkloadResourceModel maps the resource schema.
type ServerWorkloadResourceModel struct {
	// ID is required for Framework acceptance testing
	ID              types.String          `tfsdk:"id"`
	Name            types.String          `tfsdk:"name"`
	Description     types.String          `tfsdk:"description"`
	IsActive        types.Bool            `tfsdk:"is_active"`
	Tags            types.Map             `tfsdk:"tags"`
	TagsAll         types.Map             `tfsdk:"tags_all"`
	ServiceEndpoint *ServiceEndpointModel `tfsdk:"service_endpoint"`
}

// serverWorkloadDataSourceModel maps the datasource schema.
type ServerWorkloadsDataSourceModel struct {
	ServerWorkloads []ServerWorkloadResourceModel `tfsdk:"server_workloads"`
}

// ServiceEndpointModel maps service endpoint data.
type ServiceEndpointModel struct {
	ExternalID        types.String `tfsdk:"external_id"`
	ID                types.Int64  `tfsdk:"id"`
	Host              types.String `tfsdk:"host"`
	AppProtocol       types.String `tfsdk:"app_protocol"`
	TransportProtocol types.String `tfsdk:"transport_protocol"`
	RequestedPort     types.Int64  `tfsdk:"requested_port"`
	RequestedTLS      types.Bool   `tfsdk:"requested_tls"`
	Port              types.Int64  `tfsdk:"port"`
	TLS               types.Bool   `tfsdk:"tls"`

	WorkloadServiceAuthentication *WorkloadServiceAuthenticationModel `tfsdk:"authentication_config"`
	TLSVerification               types.String                        `tfsdk:"tls_verification"`
	HTTPHeaders                   types.Map                           `tfsdk:"http_headers"`
}

// WorkloadServiceAuthenticationModel maps the WorkloadServiceAuthenticationDTO struct.
type WorkloadServiceAuthenticationModel struct {
	Method types.String `tfsdk:"method"`
	Scheme types.String `tfsdk:"scheme"`
	Config types.String `tfsdk:"config"`
}
