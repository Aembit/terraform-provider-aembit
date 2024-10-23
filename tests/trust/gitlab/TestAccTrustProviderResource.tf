provider "aembit" {
}

resource "aembit_trust_provider" "gitlab1" {
	name = "TF Acceptance GitLab Job1"
	is_active = true
	gitlab_job = {
		namespace_path = "namespace_path"
		project_path = "project_path"
		ref_path = "ref_path"
		subject = "subject1"
	}
}

resource "aembit_trust_provider" "gitlab2" {
	name = "TF Acceptance GitLab Job2"
	is_active = true
	gitlab_job = {
		namespace_paths = ["namespace_path1","namespace_path2"]
		project_paths = ["project_path1","project_path2"]
		ref_paths = ["ref_path1","ref_path2"]
		subjects = ["subject1","subject2"]
	}
}

resource "aembit_trust_provider" "gitlab_mixed" {
	name = "TF Acceptance GitLab Mixed"
	is_active = true
	gitlab_job = {
		namespace_path = "namespace_path1"
		ref_path = "ref_path1"
		subjects = ["subject1","subject2"]
	}
}