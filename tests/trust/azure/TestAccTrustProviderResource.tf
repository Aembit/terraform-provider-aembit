provider "aembit" {
}

resource "aembit_trust_provider" "azure" {
	name = "TF Acceptance Azure"
	is_active = false
	azure_metadata = {
		subscription_id = "subscription_id"
	}
}