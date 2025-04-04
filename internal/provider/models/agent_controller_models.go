package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.AgentControllerResourceModel maps the resource schema.
type AgentControllerResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	IsActive           types.Bool   `tfsdk:"is_active"`
	Tags               types.Map    `tfsdk:"tags"`
	TrustProviderID    types.String `tfsdk:"trust_provider_id"`
	AllowedTLSHostname types.String `tfsdk:"allowed_tls_hostname"`
}

// agentControllerDataSourceModel maps the datasource schema.
type AgentControllersDataSourceModel struct {
	AgentControllers []AgentControllerResourceModel `tfsdk:"agent_controllers"`
}
