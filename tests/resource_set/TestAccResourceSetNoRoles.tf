provider "aembit" {
}

resource "aembit_resource_set" "no_roles" {
  name        = "TF Acceptance ResourceSet No Roles"
  description = "ResourceSet created without specifying roles"
}
