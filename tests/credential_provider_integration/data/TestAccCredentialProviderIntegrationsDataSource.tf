provider "aembit" {
}

resource "aembit_credential_provider_integration" "gitlab" {
	name = "TF Acceptance GitLab Credential Integration"
	description = "TF Acceptance GitLab Credential Integration"
	gitlab = {
		url = "https://url.com"
		personal_access_token = "test"
	}
}

data "aembit_credential_provider_integrations" "test" {
    depends_on = [ aembit_credential_provider_integration.gitlab ]
}
