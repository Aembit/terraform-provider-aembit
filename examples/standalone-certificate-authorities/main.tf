terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_standalone_certificate_authorities" "first" {}

output "first_standalone_certificate_authorities" {
  value = data.aembit_standalone_certificate_authorities.first
}
