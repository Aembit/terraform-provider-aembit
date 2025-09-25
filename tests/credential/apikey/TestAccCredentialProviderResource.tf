provider "aembit" {
  default_tags {
    tags = {
      Name           = "Terraform Agent"
      Owner          = "Aembit DevOps"
    }
  }
}

resource "aembit_credential_provider" "api_key" {
	name = "TF Acceptance API Key"
	is_active = true
	api_key = {
		api_key = "test_api_key"
	}
}
