provider "aembit" {
}

resource "aembit_resource_set" "empty_roles" {
  name        = "TF Acceptance ResourceSet Empty Roles"
  description = "ResourceSet created with explicit empty roles"
  roles       = []
}
