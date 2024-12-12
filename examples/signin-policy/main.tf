terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_signin_policy" "settings" {
  sso_required = false
  mfa_required = false
}
