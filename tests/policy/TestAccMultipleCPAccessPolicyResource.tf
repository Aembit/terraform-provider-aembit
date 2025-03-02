provider "aembit" {
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


resource "aembit_credential_provider" "snowflake1" {
	name = "TF Acceptance Snowflake Token 1"
	is_active = true
	snowflake_jwt = {
		account_id = "account-id"
		username = "username"
	}
}

resource "aembit_credential_provider" "snowflake2" {
	name = "TF Acceptance Snowflake Token 2"
	is_active = true
	snowflake_jwt = {
		account_id = "account-id"
		username = "username"
	}
}

resource "aembit_credential_provider" "snowflake3" {
	name = "TF Acceptance Snowflake Token 2"
	is_active = true
	snowflake_jwt = {
		account_id = "account-id"
		username = "username"
	}
}

resource "aembit_access_policy" "multi_cp_first_policy" {
    is_active = false
    name = "TF Multi CP First Policy"
    client_workload = aembit_client_workload.first_client.id
    credential_providers = [{
		credential_provider_id = aembit_credential_provider.snowflake1.id,
		mapping_type = "HttpHeader",
        header_name = "test_header_name",
        header_value = "test_header_value"
	}, {
		credential_provider_id = aembit_credential_provider.snowflake2.id,
        mapping_type = "HttpBody",
        httpbody_field_path = "test_field_path",
        httpbody_field_value = "test_field_value"
	}]
    server_workload = aembit_server_workload.first_server.id
}

resource "aembit_access_policy" "multi_cp_second_policy" {
    is_active = false
    name = "TF Multi CP Second Policy"
    client_workload = aembit_client_workload.second_client.id
    credential_providers = [{
		credential_provider_id = aembit_credential_provider.snowflake1.id,
		mapping_type = "AccountName",
		account_name = "account_name_1"
	}, {
		credential_provider_id = aembit_credential_provider.snowflake2.id,
        mapping_type = "AccountName",
        account_name = "account_name_2"
	}, {
		credential_provider_id = aembit_credential_provider.snowflake3.id,
        mapping_type = "AccountName",
        account_name = "account_name_3"
	}]
    server_workload = aembit_server_workload.snowflake_server.id
}
