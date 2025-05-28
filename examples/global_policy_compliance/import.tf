import {
  to = aembit_global_policy_compliance.example
  # Aembit Terraform Provider maintains a single copy of Global Policy Compliance
  # settings. Terraform requires the 'id' parameter to be present in the import {} block,
  # but its value is ignored by the Aembit Terraform Provider.
  id = "not_important"
}