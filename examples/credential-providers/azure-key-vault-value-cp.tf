resource "aembit_credential_provider" "azure_key_vault_value_cp" {
  name        = "Azure Key Vault Value Credential Provider"
  description = "Detailed description of Azure Key Vault Value Credential Provider"
  is_active   = true
  azure_key_vault_value = {
    secret_name_1                      = "secret1"
    secret_name_2                      = "secret2"
    private_network_access             = false
    credential_provider_integration_id = aembit_credential_provider_integration.azure_entra_federation_cpi.id
  }
}