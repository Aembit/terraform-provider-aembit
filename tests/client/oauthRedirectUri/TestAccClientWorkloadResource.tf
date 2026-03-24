provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

// Create a Provider and Resource in the second Resource Set
provider "aembit" {
    alias = "rs_manager"
    resource_set_id = data.aembit_resource_sets.all.resource_sets[1].id
}

resource "aembit_identity_provider" "first" {
    provider = aembit.rs_manager

    name = "TF Acceptance Client Workload First IDP"
    description = "Acceptance Test Identity Provider"
    is_active = true

    oidc = {
        oidc_base_url = "https://first.oidc.example.com"
        client_id = "first_client_id"
        scopes = "openid profile email"
        auth_type = "ClientSecret"
        client_secret = "first_secret"
        pcke_required = true
    }
}

resource "aembit_identity_provider" "second" {
    provider = aembit.rs_manager

    name = "TF Acceptance Client Workload Second IDP"
    description = "Acceptance Test Identity Provider"
    is_active = true

    oidc = {
        oidc_base_url = "https://second.oidc.example.com"
        client_id = "second_client_id"
        scopes = "openid profile email"
        auth_type = "ClientSecret"
        client_secret = "second_secret"
        pcke_required = true
    }
}

resource "aembit_client_workload" "test" {
    provider = aembit.rs_manager

    name = "TF Acceptance - Oauth Redirect URI"
    description = "Acceptance Test Client Workload"
    is_active = false
    enforce_sso = true
    identities = [
        {
            type = "oauthRedirectUri"
            value = "https://test.aembit.local:12345"
        }
    ]
    sso_identity_providers = [
        aembit_identity_provider.first.id,
        aembit_identity_provider.second.id,
    ]
}
