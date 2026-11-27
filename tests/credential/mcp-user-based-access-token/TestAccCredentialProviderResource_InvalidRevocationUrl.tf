provider "aembit" {
}

resource "aembit_credential_provider" "mcp_user_based_access_token" {
	name = "TF Acceptance McpUserBasedAccessToken Invalid Revocation URL"
	is_active = true
	mcp_user_based_access_token = {
		mcp_server_url = "https://aembit.io/.well-known/openid-configuration"
		oauth_authorization_url = "https://aembit.io/authorize"
		oauth_token_url = "https://aembit.io/token"
		oauth_revocation_url = "http://aembit.io/revoke"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		is_pkce_required = true
	}
}
