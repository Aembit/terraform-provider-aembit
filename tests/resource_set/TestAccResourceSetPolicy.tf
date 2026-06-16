provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_resource_set" "crs" {
	name = "TF Acceptance Custom Policy ResourceSet"
	description = "TF Acceptance Custom Policy ResourceSet"
	roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
}

resource "aembit_client_workload" "test" {
	resource_set_id = aembit_resource_set.crs.id
    name = "Unit Test 1"
    description = "Acceptance Test client workload"
    is_active = true
    identities = [
        {
            type = "k8sNamespace"
            value = "unittest1namespace"
        },
    ]
}

resource "aembit_server_workload" "test" {
	resource_set_id = aembit_resource_set.crs.id
	name = "Unit Test 1"
    description = "Description"
    is_active = true
	service_endpoint = {
		host = "unittest.testhost.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			method = "HTTP Authentication"
			scheme = "Bearer"
		}
		http_headers = {
			host = "graph.microsoft.com"
			user-agent = "curl/7.64.1"
			accept = "*/*"
		}
	}
}

resource "aembit_integration" "wiz" {
    resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://api.us17.app.wiz.io/graphql"
	oauth_client_credentials = {
		token_url = "https://auth.app.wiz.io/oauth/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "test" {
    resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance Wiz"
	is_active = false
	integration_id = aembit_integration.wiz.id
	wiz_conditions = {
		max_last_seen = 3600
		container_cluster_connected = true
	}
}

resource "aembit_credential_provider_integration" "azure_entra_federation_cpi" {
	resource_set_id = aembit_resource_set.crs.id    
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

resource "aembit_credential_provider" "test" {
	resource_set_id = aembit_resource_set.crs.id    
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

resource "aembit_trust_provider" "test" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance GitLab Job1"
	is_active = true
	gitlab_job = {
		namespace_path = "namespace_path"
		project_path = "project_path"
		ref_path = "ref_path"
		subject = "subject1"
	}
}

resource "aembit_standalone_certificate_authority" "test" {
	resource_set_id = aembit_resource_set.crs.id
	name = "unittestname"
    description = "Description"
	leaf_lifetime = 1440
}

resource "aembit_access_policy" "first_policy" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF ResourceSet Policy"
    is_active = true
    client_workload = aembit_client_workload.test.id
    trust_providers = [
        aembit_trust_provider.test.id
    ]
    access_conditions = [
        aembit_access_condition.test.id
    ]
    credential_provider = aembit_credential_provider.test.id
    server_workload = aembit_server_workload.test.id
}

