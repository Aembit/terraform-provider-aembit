provider "aembit" {
}

resource "aembit_server_workload" "snowflake_server" {
    name = "snowflake terraform server workload"
    description = "new snowflake server workload for policy integration"
    is_active = false
    service_endpoint = {
        host = "myhost.unittest.com"
        port = 443
        app_protocol = "Snowflake"
		transport_protocol = "TCP"
        requested_port = 80
        tls_verification = "full"
	    requested_tls = true
	    tls = true
    }
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

resource "aembit_credential_provider" "api_key" {
	name = "TF Acceptance Policy CP"
	api_key = {
		api_key = "test"
	}
}

resource "aembit_access_policy" "multi_cp_duplicate_policy_2" {
    is_active = false
    name = "TF Multi CP Duplicate Policy"
    client_workload = aembit_client_workload.first_client.id
    credential_provider = aembit_credential_provider.api_key.id
    server_workload = aembit_server_workload.snowflake_server.id
}
