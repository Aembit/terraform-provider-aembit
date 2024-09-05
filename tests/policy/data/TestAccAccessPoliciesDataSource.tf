provider "aembit" {
}

resource "aembit_client_workload" "first_client" {
    name = "first terraform client workload"
    description = "new client workload for policy integration"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "clientworkloadNamespace"
        },
    ]
}

resource "aembit_trust_provider" "azure1" {
	name = "TF Acceptance Azure"
	azure_metadata = {
		subscription_id = "subscription_id"
	}
}

resource "aembit_trust_provider" "azure2" {
	name = "TF Acceptance Azure"
	azure_metadata = {
		subscription_id = "subscription_id"
	}
}

resource "aembit_integration" "wiz" {
	name = "TF Acceptance Wiz"
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://endpoint"
	oauth_client_credentials = {
		token_url = "https://url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "wiz" {
	name = "TF Acceptance Wiz"
	integration_id = aembit_integration.wiz.id
	wiz_conditions = {
		max_last_seen = 3600
		container_cluster_connected = true
	}
}

resource "aembit_credential_provider" "snowflake1" {
	name = "TF Acceptance Snowflake Token 1"
	is_active = true
	snowflake_jwt = {
		account_id = "account_id"
		username = "username"
	}
}

resource "aembit_server_workload" "first_server" {
    name = "first terraform server workload"
    description = "new server workload for policy integration"
    is_active = false
    service_endpoint = {
        host = "myhost.unittest.com"
        port = 443
        app_protocol = "HTTP"
		transport_protocol = "TCP"
        requested_port = 80
        tls_verification = "full"
	    requested_tls = true
	    tls = true
    }
}

resource "aembit_access_policy" "first_policy" {
    is_active = false
    client_workload = aembit_client_workload.first_client.id
    trust_providers = [
        aembit_trust_provider.azure1.id,
        aembit_trust_provider.azure2.id
    ]
    access_conditions = [
        aembit_access_condition.wiz.id
    ]
    credential_providers = [{
		credential_provider_id = aembit_credential_provider.snowflake1.id,
		mapping_type = "None",
        header_name = "",
        header_value = "",
		account_name = "",
		httpbody_field_path = "",
        httpbody_field_value = ""
	}]
    server_workload = aembit_server_workload.first_server.id
}

data "aembit_access_policies" "test" {}
