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

resource "aembit_credential_provider" "gitlab_managed_account" {
	name = "TF Acceptance Managed Gitlab Account"
	is_active = true
	managed_gitlab_account = {
		group_ids = ["678","test_group_id"]
		project_ids = ["123","456"]
		access_level = 30
		lifetime_in_days = 2
		scope = "api test"
		credential_provider_integration_id = aembit_credential_provider_integration.gitlab.id
	}
}
