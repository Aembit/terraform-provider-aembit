provider "aembit" {
}

variable "integrations" {
   type = map(object({
    name = string
    is_active = bool
	  type = string
	  sync_frequency = number
	  endpoint = string
    oauth_client_credentials = map(string)
  }))
  default = {
    wiz = {
      name = "TF Acceptance Wiz"
      is_active = true
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
    crowdstrike = {
      name = "TF Acceptance Crowdstrike"
	    is_active = true
	    type = "CrowdStrike"
	    sync_frequency = 3600
	    endpoint = "https://endpoint"
      oauth_client_credentials = {
        token_url = "https://url/token"
        client_id = "client_id"
        client_secret = "client_secret"
        audience = "audience"
      }
    }
  }
}

variable "access_conditions" {
  type = map(object({
    max_last_seen               = number
    match_hostname              = optional(bool, null)
    match_serial_number         = optional(bool, null)
    prevent_rfm                 = optional(bool, null)
    container_cluster_connected = optional(bool, null)
  }))
  default = {
    crowdstrike = {
      max_last_seen               = 3600
      match_hostname              = true
      match_serial_number         = true
      prevent_rfm                 = true
    }
    wiz = {
      max_last_seen               = 3600
      container_cluster_connected = true
    }
  }
}

resource "aembit_integration" "ac" {
  for_each                  = var.integrations
  name                      = each.value.name
  is_active                 = each.value.is_active
  type                      = each.value.type
  sync_frequency            = each.value.sync_frequency
  endpoint                  = each.value.endpoint
  oauth_client_credentials  = each.value.oauth_client_credentials
}

resource "aembit_access_condition" "ac" {
  for_each        = var.access_conditions
  name            = "ac-${each.key}"
  is_active       = each.value.is_active
  integration_id  = aembit_integration.ac[each.key].id
  
  dynamic "crowdstrike_conditions" {
    count = var.integration_type == "CrowdStrike" ? 1 : 0
    content {
      max_last_seen       = each.value.max_last_seen
      match_hostname      = each.value.match_hostname
      match_serial_number = each.value.match_serial_number
      prevent_rfm         = each.value.prevent_rfm
    }
  }

  dynamic "wiz_conditions" {
    count = var.integration_type == "WizIntegrationApi" ? 1 : 0
    content {
      max_last_seen               = each.value.max_last_seen
      container_cluster_connected = each.value.container_cluster_connected
    }
  }
}


