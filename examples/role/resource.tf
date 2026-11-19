# 1. Default Resource Set assignment: role created with resource_sets_assignments omitted
# Automatically defaults to the tenant's Default Resource Set ("ffffffff-ffff-ffff-ffff-ffffffffffff")
resource "aembit_role" "default_resource_set" {
  name        = "Security Operator - Default Scope"
  description = "Role scoped to the Default Resource Set"
  is_active   = true

  access_policies                    = { read = true, write = false }
  client_workloads                   = { read = true, write = false }
  trust_providers                    = { read = true, write = false }
  access_conditions                  = { read = true, write = false }
  integrations                       = { read = true, write = false }
  credential_providers               = { read = true, write = false }
  credential_provider_integrations   = { read = true, write = false }
  server_workloads                   = { read = true, write = false }
  agent_controllers                  = { read = true, write = false }
  discovery_integrations             = { read = true, write = false }
  standalone_certificate_authorities = { read = true, write = false }

  access_authorization_events     = { read = true }
  global_policy_compliance_report = { read = true }
  audit_logs                      = { read = true }
  workload_events                 = { read = true }

  users                    = { read = true, write = false }
  roles                    = { read = true, write = false }
  log_streams              = { read = true, write = false }
  identity_providers       = { read = true, write = false }
  signon_policy            = { read = true, write = false }
  global_policy_compliance = { read = true, write = false }
  resource_sets            = { read = true, write = false }
  routing                  = { read = true, write = false }

  # resource_sets_assignments omitted -> defaults to ["ffffffff-ffff-ffff-ffff-ffffffffffff"]
}

# 2. Scoped to a custom Resource Set (Pattern A - Recommended):
# aembit_role owns the association via resource_sets_assignments.
# Prerequisite: resource_sets = { read = true, ... } is REQUIRED when assigning non-Default Resource Sets.
# Note: roles is omitted on aembit_resource_set.custom_app to avoid conflicting bidirectional state.
resource "aembit_resource_set" "custom_app" {
  name        = "App Team Production"
  description = "Resource Set for production application workloads"
}

resource "aembit_role" "custom_resource_set" {
  name        = "App Operator - Custom Scope"
  description = "Role scoped to a custom Resource Set"
  is_active   = true

  access_policies                    = { read = true, write = true }
  client_workloads                   = { read = true, write = true }
  trust_providers                    = { read = true, write = true }
  access_conditions                  = { read = true, write = true }
  integrations                       = { read = true, write = true }
  credential_providers               = { read = true, write = true }
  credential_provider_integrations   = { read = true, write = true }
  server_workloads                   = { read = true, write = true }
  agent_controllers                  = { read = true, write = true }
  discovery_integrations             = { read = true, write = true }
  standalone_certificate_authorities = { read = true, write = true }

  access_authorization_events     = { read = true }
  global_policy_compliance_report = { read = true }
  audit_logs                      = { read = true }
  workload_events                 = { read = true }

  users                    = { read = true, write = false }
  roles                    = { read = true, write = false }
  log_streams              = { read = true, write = false }
  identity_providers       = { read = true, write = false }
  signon_policy            = { read = true, write = false }
  global_policy_compliance = { read = true, write = false }
  resource_sets            = { read = true, write = false } # Required: read = true for custom Resource Sets
  routing                  = { read = true, write = false }

  resource_sets_assignments = [aembit_resource_set.custom_app.id]
}

# 3. Association managed from aembit_resource_set.roles (Pattern B):
# aembit_resource_set owns the association. Set resource_sets_assignments = []
# and ignore changes on resource_sets_assignments to prevent default assignment drift.
resource "aembit_role" "rs_managed_role" {
  name        = "Data Analyst - ResourceSet Managed"
  description = "Role whose Resource Set association is managed from aembit_resource_set.roles"
  is_active   = true

  access_policies                    = { read = true, write = false }
  client_workloads                   = { read = true, write = false }
  trust_providers                    = { read = true, write = false }
  access_conditions                  = { read = true, write = false }
  integrations                       = { read = true, write = false }
  credential_providers               = { read = true, write = false }
  credential_provider_integrations   = { read = true, write = false }
  server_workloads                   = { read = true, write = false }
  agent_controllers                  = { read = true, write = false }
  discovery_integrations             = { read = true, write = false }
  standalone_certificate_authorities = { read = true, write = false }

  access_authorization_events     = { read = true }
  global_policy_compliance_report = { read = true }
  audit_logs                      = { read = true }
  workload_events                 = { read = true }

  users                    = { read = true, write = false }
  roles                    = { read = true, write = false }
  log_streams              = { read = true, write = false }
  identity_providers       = { read = true, write = false }
  signon_policy            = { read = true, write = false }
  global_policy_compliance = { read = true, write = false }
  resource_sets            = { read = true, write = false }
  routing                  = { read = true, write = false }

  # Opt out of default Default Resource Set assignment and ignore out-of-band updates
  resource_sets_assignments = []
  lifecycle {
    ignore_changes = [resource_sets_assignments]
  }
}

resource "aembit_resource_set" "analytics" {
  name        = "Analytics Platform"
  description = "Resource Set managing associated roles"
  roles       = [aembit_role.rs_managed_role.id]
}
