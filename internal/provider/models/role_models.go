package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleResourceModel maps the resource schema.
type RoleResourceModel struct {
	// ID is required for Framework acceptance testing
	ID             types.String    `tfsdk:"id"`
	Name           types.String    `tfsdk:"name"`
	Description    types.String    `tfsdk:"description"`
	IsActive       types.Bool      `tfsdk:"is_active"`
	Tags           types.Map       `tfsdk:"tags"`
	AccessPolicies *RolePermission `tfsdk:"access_policies"`
	Routing        *RolePermission `tfsdk:"routing"`

	ClientWorkloads     *RolePermission `tfsdk:"client_workloads"`
	TrustProviders      *RolePermission `tfsdk:"trust_providers"`
	AccessConditions    *RolePermission `tfsdk:"access_conditions"`
	Integrations        *RolePermission `tfsdk:"integrations"`
	CredentialProviders *RolePermission `tfsdk:"credential_providers"`
	ServerWorkloads     *RolePermission `tfsdk:"server_workloads"`

	AgentControllers                 *RolePermission `tfsdk:"agent_controllers"`
	StandaloneCertificateAuthorities *RolePermission `tfsdk:"standalone_certificate_authorities"`

	AccessAuthorizationEvents *RoleReadOnlyPermission `tfsdk:"access_authorization_events"`
	AuditLogs                 *RoleReadOnlyPermission `tfsdk:"audit_logs"`
	WorkloadEvents            *RoleReadOnlyPermission `tfsdk:"workload_events"`

	Users             *RolePermission `tfsdk:"users"`
	SignOnPolicy      *RolePermission `tfsdk:"signon_policy"`
	Roles             *RolePermission `tfsdk:"roles"`
	ResourceSets      *RolePermission `tfsdk:"resource_sets"`
	LogStreams        *RolePermission `tfsdk:"log_streams"`
	IdentityProviders *RolePermission `tfsdk:"identity_providers"`
}

// roleDataSourceModel maps the datasource schema.
type RolesDataSourceModel struct {
	Roles []RoleResourceModel `tfsdk:"roles"`
}

type RolePermission struct {
	Read  types.Bool `tfsdk:"read"`
	Write types.Bool `tfsdk:"write"`
}

type RoleReadOnlyPermission struct {
	Read types.Bool `tfsdk:"read"`
}
