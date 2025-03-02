provider "aembit" {
}


resource "aembit_role" "role" {
	name = "TF Acceptance Role for Token"
	is_active = true
	access_policies = { read = true, write = true }
	routing = { read = true, write = true }
	client_workloads = { read = true, write = true }
	trust_providers = { read = true, write = true }
	access_conditions = { read = true, write = true }
	integrations = { read = true, write = true }
	credential_providers = { read = true, write = true }
	server_workloads = { read = true, write = true }

	agent_controllers = { read = true, write = true }

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

resource "aembit_credential_provider" "aembit" {
	name = "TF Acceptance Aembit Token"
	is_active = false
	aembit_access_token = {
		role_id = aembit_role.role.id
		lifetime = 1800
	}
}
