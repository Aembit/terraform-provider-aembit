provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

locals {
    tf_testing_rs_id = [for rs in data.aembit_resource_sets.all.resource_sets : rs.id if rs.name == "TF Testing"][0]
}

provider "aembit" {
    alias = "rs_manager"
    resource_set_id = local.tf_testing_rs_id
}

resource "aembit_identity_provider" "test_idp" {
    provider = aembit.rs_manager

    name        = "TF Acceptance CW UserGroup IDP"
    description = "Acceptance Test Identity Provider for UserGroup CW"
    is_active   = true

    oidc = {
        oidc_base_url = "https://usergroup.oidc.example.com"
        client_id     = "usergroup_client_id"
        scopes        = "openid profile email"
        auth_type     = "ClientSecret"
        client_secret = "usergroup_secret"
        pcke_required = true
    }
}

resource "aembit_client_workload" "test" {
    provider = aembit.rs_manager

    name        = "TF Acceptance - UserGroup"
    description = "Acceptance Test Client Workload"
    is_active   = false
    identities = [
        {
            type                 = "userGroup"
            value                = "*group_name*"
            identity_provider_id = aembit_identity_provider.test_idp.id
        }
    ]
}
