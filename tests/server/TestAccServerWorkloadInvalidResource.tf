provider "aembit" {
}

resource "aembit_server_workload" "test_mcp_invalid" {
	name = "MCP Invalid Test"
	description = "Testing MCP with mismatched ports"
	is_active = false
	service_endpoint = {
		host = "mcp.testhost.com"
		port = 443
		requested_port = 8080
		tls = true
		app_protocol = "MCP"
		url_path = "/token"
		transport_protocol = "TCP"
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
	}
}