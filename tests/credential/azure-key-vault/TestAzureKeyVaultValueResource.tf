resource "aembit_credential_provider_integration" "azure_entra_federation_cpi" {
    name                   = "TF Acceptance Azure Entra Federation Credential Provider Integration"
    description            = "TF Acceptance Azure Entra Federation Credential Provider Integration Description"
    azure_entra_federation = {
        audience       = "api://AzureADTokenExchange"
        subject        = "subject"
        azure_tenant   = "00000000-0000-0000-0000-000000000000"
        client_id      = "00000000-0000-0000-0000-000000000000"
        key_vault_name = "KeyVaultName"
    }
}

resource "aembit_credential_provider" "azure_key_vault_value_cp" {
    name                  = "TF Acceptance Azure Key Vault Value CP"
    description           = "TF Acceptance Azure Key Vault Value CP Description"
    is_active             = true
    azure_key_vault_value = {
        secret_name_1                      = "secret1"
        secret_name_2                      = "secret2"
        private_network_access             = false
        credential_provider_integration_id = aembit_credential_provider_integration.azure_entra_federation_cpi.id
    }
}