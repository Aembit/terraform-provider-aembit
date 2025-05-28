resource "aembit_global_policy_compliance" "example" {
  access_policy_trust_provider_compliance          = "Recommended"
  access_policy_access_condition_compliance        = "Recommended"
  agent_controller_trust_provider_compliance       = "Recommended"
  agent_controller_allowed_tls_hostname_compliance = "Recommended"
}