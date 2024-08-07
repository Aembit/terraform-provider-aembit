terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_credential_providers" "first" {
}

output "first" {
  value     = data.aembit_credential_providers.first
  sensitive = true
}
