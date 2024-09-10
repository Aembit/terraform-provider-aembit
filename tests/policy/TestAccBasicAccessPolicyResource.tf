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

resource "aembit_client_workload" "second_client" {
    name = "second terraform client workload"
    description = "new client workload for policy integration"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "secondClientWorkloadNamespace"
        },
    ]
}

resource "aembit_client_workload" "third_client" {
    name = "third terraform client workload"
    description = "new client workload for policy integration"
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = "thirdClientWorkloadNamespace"
        },
    ]
}

resource "aembit_credential_provider" "snowflake1" {
	name = "TF Acceptance Snowflake Token 1"
	is_active = true
	snowflake_jwt = {
		account_id = "account_id"
		username = "username"
	}
}

resource "aembit_credential_provider" "snowflake2" {
	name = "TF Acceptance Snowflake Token 2"
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
    name = "TF First Policy"
    client_workload = aembit_client_workload.first_client.id
    trust_providers = []
    access_conditions = []
    credential_provider = aembit_credential_provider.snowflake1.id
    server_workload = aembit_server_workload.first_server.id
}

resource "aembit_access_policy" "second_policy" {
    is_active = false
    name = "TF Second Policy"
    client_workload = aembit_client_workload.second_client.id
    trust_providers = []
    access_conditions = []
    credential_provider = aembit_credential_provider.snowflake2.id
    server_workload = aembit_server_workload.first_server.id
}

resource "aembit_access_policy" "third_policy" {
    is_active = false
    name = "TF Third Policy"
    client_workload = aembit_client_workload.third_client.id
    credential_providers = [{
		credential_provider_id = aembit_credential_provider.snowflake1.id,
		mapping_type = "HttpHeader",
        header_name = "test_header_name",
        header_value = "test_header_value",
		account_name = "",
		httpbody_field_path = "",
        httpbody_field_value = ""
	},
	{
		credential_provider_id = aembit_credential_provider.snowflake2.id,
        mapping_type = "HttpBody",
        header_name = "",
        header_value = "",
        account_name = "",
        httpbody_field_path = "test_field_path",
        httpbody_field_value = "test_field_value"
	}]
    server_workload = aembit_server_workload.first_server.id
}

