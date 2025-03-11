provider "aembit" {
}

resource "aembit_role" "role" {
	name = "TF Acceptance Role"
	is_active = false
	access_policies = { read = true, write = true }
	routing = { read = true, write = true }
	client_workloads = { read = true, write = true }
	trust_providers = { read = true, write = true }
	access_conditions = { read = true, write = true }
	integrations = { read = true, write = true }
	credential_providers = { read = true, write = true }
	server_workloads = { read = true, write = true }

	agent_controllers = { read = true, write = true }
	standalone_certificate_authorities = { read = true, write = true }

	access_authorization_events = { read = true }
	audit_logs = { read = true }
	workload_events = { read = true }

	users = { read = true, write = true }
	signon_policy = { read = true, write = true }
	roles = { read = true, write = true }
	resource_sets = { read = true, write = true }
	log_streams = { read = true, write = true }
	identity_providers = { read = true, write = true }
}