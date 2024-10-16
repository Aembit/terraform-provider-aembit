provider "aembit" {
}

resource "aembit_trust_provider" "gitlab" {
	name = "TF Acceptance GitLab Job"
	is_active = true
	gitlab_job = {
		namespace_path = "namespace_path"
		project_path = "project_path"
		ref_path = "ref_path"
		subject = "subject"
	}
}