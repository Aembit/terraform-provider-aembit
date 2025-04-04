provider "aembit" {
}

resource "aembit_agent_controller" "device_code" {
	name = "TF Acceptance Device Code"
    description = "device code agent controller"
	is_active = true
    
}

resource "aembit_agent_controller" "azure_tp" {
	name = "TF Acceptance Azure Trust Provider"
	is_active = false
    tags = {
        color = "blue"
        day   = "Sunday"
    }

	trust_provider_id = aembit_trust_provider.azure.id
	allowed_tls_hostname = "test.example.com"
}

resource "aembit_trust_provider" "azure" {
	name = "TF Acceptance Azure"
	azure_metadata = {
		subscription_id = "subscription_id"
	}
}