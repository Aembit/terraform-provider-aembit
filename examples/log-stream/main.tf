terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

data "aembit_log_streams" "first" {
}

output "first" {
  value     = data.aembit_log_streams.first
  sensitive = true
}
