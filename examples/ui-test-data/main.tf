terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
}

resource "aembit_client_workload" "cw" {
    for_each = { for i in range(1, 11)  : i => {
        name  = "ClientWorkload UI Test Lorem ipsum dolor sit amet. consectetur adipiscing elit ${i}"
        identity_namespace = "namespace${i}value"
        identity_podname = "podname${i}value"
        identity_podname_prefix = "podnameprefix${i}value"
        identity_service_accountname = "accountname${i}value"
        description = "Etiam vel arcu tristique. feugiat odio at. sagittis nunc. Cras at tristique ex. ${i}"
    } }

    name        = each.value.name
    description = each.value.description
    is_active = false
    identities = [
        {
            type = "k8sNamespace"
            value = each.value.identity_namespace
        },
        {
            type = "k8sPodName"
            value = each.value.identity_podname
        },
        {
            type = "k8sPodNamePrefix"
            value = each.value.identity_podname_prefix
        },
        {
            type = "k8sServiceAccountName"
            value = each.value.identity_service_accountname
        },
    ]
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

resource "aembit_trust_provider" "tp" {
    for_each = { for i in range(1, 11)  : i => {
        name  = "TrustProvider UI Test Lorem ipsum dolor sit amet. consectetur adipiscing elit ${i}"
        description = "Etiam vel arcu tristique. feugiat odio at. sagittis nunc. Cras at tristique ex. ${i}"
    } }

    name        = each.value.name
    description = each.value.description
    is_active = false
	aws_role = {
		account_id = "account_id"
		assumed_role = "assumed_role"
		role_arn = "role_arn"
		username = "username"
	}
}

resource "aembit_credential_provider" "cp" {
    for_each = { for i in range(1, 11)  : i => {
        name  = "CredentialProvider UI Test Lorem ipsum dolor sit amet. consectetur adipiscing elit ${i}"
        description = "Etiam vel arcu tristique. feugiat odio at. sagittis nunc. Cras at tristique ex. ${i}"
        api_key = "test_api_key${i}"
    } }

    name        = each.value.name
    description = each.value.description
    is_active = false
	api_key = {
		api_key = each.value.api_key
	}
}

resource "aembit_server_workload" "sw" {
    for_each = { for i in range(1, 11)  : i => {
        name  = "ServerWorkload UI Test Lorem ipsum dolor sit amet. consectetur adipiscing elit ${i}"
        description = "Etiam vel arcu tristique. feugiat odio at. sagittis nunc. Cras at tristique ex. ${i}"
        api_key = "test_api_key${i}"
    } }

    name        = each.value.name
    description = each.value.description
    is_active = false
	service_endpoint = {
		host = "uitest.testhost.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			"method" = "HTTP Authentication"
			"scheme" = "Bearer"
		}
		http_headers = {
			host = "graph.microsoft.com"
			user-agent = "curl/7.64.1"
			accept = "*/*"
		}
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

resource "aembit_integration" "wiz" {
    for_each = { for i in range(1, 11)  : i => {
        name  = "Integration UI Test Lorem ipsum dolor sit amet. consectetur adipiscing elit ${i}"
        description = "Etiam vel arcu tristique. feugiat odio at. sagittis nunc. Cras at tristique ex. ${i}"
    } }

    name        = each.value.name
    description = each.value.description
    is_active = false
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://endpoint"
	oauth_client_credentials = {
		token_url = "https://url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "ac" {
    for_each = {
        for key, value in aembit_integration.wiz : key => {
            integration_id = aembit_integration.wiz[key].id
            integration_name = value.name
            ac_name = "AccessCondition UI Test Lorem ipsum dolor sit amet. consectetur adipiscing"
        }
    }

	name = each.value.ac_name
	is_active = false
	integration_id =  each.value.integration_id
	wiz_conditions = {
		max_last_seen = 3600
		container_cluster_connected = true
	}

    depends_on = [ aembit_integration.wiz ]
}

resource "aembit_access_policy" "policy" {
    for_each = {
        for key, value in aembit_client_workload.cw : key => {
            cw_id = aembit_client_workload.cw[key].id
            sw_id = aembit_server_workload.sw[key].id
            tp_id = aembit_trust_provider.tp[key].id
            cp_id = aembit_credential_provider.cp[key].id
            ac_id = aembit_access_condition.ac[key].id
        }
    }

	name        = "Placeholder"
    is_active = true
    client_workload = each.value.cw_id
    server_workload = each.value.sw_id
    trust_providers = [
        each.value.tp_id
    ]
    access_conditions = [
        each.value.ac_id
    ]
    credential_provider =  each.value.cp_id
}
