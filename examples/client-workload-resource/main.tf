terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_identity_provider" "primary" {
  name        = "terraform example primary idp"
  description = "example identity provider"
  is_active   = true

  oidc = {
    oidc_base_url = "https://example.oidc.com"
    client_id     = "example-client-id"
    scopes        = "openid profile email"
    auth_type     = "ClientSecret"
    client_secret = "example-client-secret"
    pcke_required = true
  }
}

resource "aembit_client_workload" "edu" {
  name        = "terraform client workload3"
  description = "new client workload3"
  is_active   = false
  enforce_sso = false
  identities = [
    {
      type  = "oauthRedirectUri"
      value = "https://example.com/callback"
    },
  ]
  sso_identity_providers = [aembit_identity_provider.primary.id]
}

output "edu_client_workload" {
  value = aembit_client_workload.edu
}
