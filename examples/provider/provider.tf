terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
  # This client_id configuration may be set here or in the AEMBIT_CLIENT_ID environment variable.
  # Note: This is a sample value and must be replaced with your Aembit Trust Provider generated value.
  client_id = "aembit:useast2:tenant:identity:github_idtoken:0bc4dbcd-e9c8-445b-ac90-28f47b8649cc"

  # Optional, defaults to the Default Resource Set
  # Note: This is a sample value and must be replaced with your generated Resource Set ID.
  resource_set_id = "d67afe77-6313-4b18-9b64-c0949b75bd1c"
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}
