provider "aembit" {
}

resource "aembit_trust_provider" "kerberos" {
	name = "TF Acceptance Kerberos"
	kerberos = {
		agent_controller_id = "fb1b63df-bfc2-4ff9-baa4-fc84bdf9a7e5"
		realm = "realm"
		principal = "principal"
		source_ip = "source_ip"
	}
}