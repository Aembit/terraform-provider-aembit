provider "aembit" {
}

resource "aembit_credential_provider_integration" "gitlab" {
	name = "TF Acceptance GitLab Credential Integration"
	description = "TF Acceptance GitLab Credential Integration"
	gitlab = {
		url = "https://gitlab.aembit-eng.com"
		personal_access_token = "REPLACE_WITH_TEST_TOKEN"
	}
}

resource "aembit_credential_provider" "gitlab_managed_account" {
	name = "TF Acceptance Managed Gitlab Account"
	is_active = true
	managed_gitlab_account = {
		service_account_username = "test_service_account"
		group_ids = ["678","test_group_id"]
		project_ids = ["123","456"]
		access_level = 30
		lifetime_in_days = 0.25
		lifetime_in_hours = 6
		scope = "api test"
		credential_provider_integration_id = aembit_credential_provider_integration.gitlab.id
	}
}
