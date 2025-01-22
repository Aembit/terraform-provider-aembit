provider "aembit" {
}

resource "aembit_credential_provider" "ae" {
	name = "TF Acceptance Azure Entra Workload"
	is_active = true
	azure_entra_workload_identity = {
		audience = "audience"
		subject = "subject"
		scope = "scope"
		azure_tenant = "00000000-0000-0000-0000-000000000000"
		client_id = "00000000-0000-0000-0000-000000000000"
	}
}
