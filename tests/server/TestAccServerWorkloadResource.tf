provider "aembit" {
}

resource "aembit_server_workload" "test" {
	name = "Unit Test 1"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "unittest.testhost.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
		http_headers = {
			host = "graph.microsoft.com"
			user-agent = "curl/7.64.1"
			accept = "*/*"
		}
	}
}

resource "aembit_server_workload" "test_wildcard" {
	name = "Unit Test 1 Wildcard"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "*.testhost.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
	}
}

resource "aembit_server_workload" "test_oauth" {
	name = "Unit Test 1 OAuth"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "oauth.testhost.com"
		port = 443
        tls = true
		app_protocol = "OAuth"
		url_path = "/token"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			method = "OAuth Client Authentication"
			scheme = "POST Body"
		}
	}
}

resource "aembit_server_workload" "test_mcp_empty_url_path" {
	name = "Unit Test 1 MCP"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "mcp.testhost.com"
		port = 443
        tls = true
		app_protocol = "MCP"
		transport_protocol = "TCP"
		requested_port = 1
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
	}
}

resource "aembit_server_workload" "test_mcp_with_url_path" {
	name = "Unit Test 1 MCP"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "mcp.testhost.com"
		port = 443
        tls = true
		app_protocol = "MCP"
		url_path = "/token"
		transport_protocol = "TCP"
		requested_port = 1
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
	}
}

resource "aembit_server_workload" "test_oracle_database" {
	name = "Unit Test 1 Oracle Database"
    description = "Description"
    is_active = false
	service_endpoint = {
		host = "oracleDatabase.testhost.com"
		port = 443
        tls = true
		app_protocol = "Oracle Database"
		transport_protocol = "TCP"
		requested_port = 1521
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			method = "Password Authentication"
			scheme = "Password"
		}
	}
}