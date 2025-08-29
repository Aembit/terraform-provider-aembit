provider "aembit" {
  alias = "secondary"
  # ...secondary provider config...
}

data "aembit_caller_identity" "secondary" {
  provider = aembit.secondary
}