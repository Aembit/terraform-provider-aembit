# 1. Default system roles: roles omitted (system automatically assigns SuperAdmin and Auditor; both appear in effective_roles)
resource "aembit_resource_set" "default_roles" {
  name        = "Workloads - Default Roles"
  description = "Resource Set with default system role assignments"
}

# 2. Explicit empty custom roles: roles = [] (or module/variable fallback try(each.value.roles, []))
# SuperAdmin and Auditor are still assigned by the backend and exposed in effective_roles
resource "aembit_resource_set" "empty_custom_roles" {
  name        = "Workloads - Explicit Empty Custom Roles"
  description = "Resource Set configured with explicit empty custom roles"
  roles       = []
}

# 3. Custom role assignment: resource set configured with explicit custom roles
# Note: When managing associations from aembit_role (Pattern A), omit roles on aembit_resource_set instead
resource "aembit_resource_set" "custom_roles" {
  name        = "Workloads - Custom Roles"
  description = "Resource Set configured with custom role assignments"
  roles       = [aembit_role.example.id]
}

resource "aembit_role" "example" {
  name        = "Example Custom Role"
  description = "Role associated with custom resource set"
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

  resource_sets_assignments = []
  lifecycle {
    ignore_changes = [resource_sets_assignments]
  }
}
