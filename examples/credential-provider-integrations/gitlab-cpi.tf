resource "aembit_credential_provider_integration" "gitlab" {
  name        = "Managed GitLab Account Credential Provider Integration"
  description = "Detailed description of Managed GitLab Account Credential Provider Integration"
  gitlab = {
    url                   = "https://url.com"
    personal_access_token = "sample_value"
  }
}