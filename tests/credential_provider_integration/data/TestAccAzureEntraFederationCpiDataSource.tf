provider "aembit" {
}

resource "aembit_credential_provider_integration" "azure_entra_federation_cpi" {
    name                   = "TF Acceptance Azure Entra Federation Credential Provider Integration"
    description            = "TF Acceptance Azure Entra Federation Credential Provider Integration Description"
    azure_entra_federation = {
        audience           = "api://AzureADTokenExchange"
        subject            = "subject"
        azure_tenant       = "00000000-0000-0000-0000-000000000000"
        client_id          = "00000000-0000-0000-0000-000000000000"
        key_vault_name     = "KeyVaultName"
        fetch_secret_names = false
    }
}

data "aembit_credential_provider_integrations" "aef_cpi_test" {
    depends_on = [ aembit_credential_provider_integration.azure_entra_federation_cpi ]
}