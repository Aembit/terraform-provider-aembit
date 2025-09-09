provider "aembit" {
}

resource "aembit_trust_provider" "github" {
	name = "TF Acceptance GitHub Action"
	is_active = true
	github_action = {
		actor = "actor"
		repository = "repository"
		workflow = "workflow"
		oidc_endpoint = "https://gitlab.com"
	}
}