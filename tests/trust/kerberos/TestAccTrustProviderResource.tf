provider "aembit" {
}

resource "aembit_trust_provider" "kerberos" {
	name = "TF Acceptance Kerberos"
	is_active = true
	kerberos = {
		agent_controller_ids = ["fb1b63df-bfc2-4ff9-baa4-fc84bdf9a7e5"]
		realm_domain = "realm_domain"
		principal = "principal"
		source_ip = "source_ip"
	}
}