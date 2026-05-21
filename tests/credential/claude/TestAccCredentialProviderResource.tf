provider "aembit" {
}

resource "aembit_credential_provider" "claude" {
	name = "TF Acceptance Claude Wif"
	is_active = true
	claude_wif = {
		federation_rule_id = "fdrl_test"
		service_account_id = "svac_test"
		organization_id   = "7f0a8276-8f3b-4c3e-89a3-558667666012"
		audience         = "aud_test"
		workspace_id      = "wrkspc_test"
		scope            = "scope_test"
		lifetime         = 3600
	}
}
