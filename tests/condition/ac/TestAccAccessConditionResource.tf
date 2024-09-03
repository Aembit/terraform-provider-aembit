provider "aembit" {}

variable "integrations" {
  type = list(object({
    name = string
    is_active = bool
    type = string
    sync_frequency = number
    endpoint = string
    token_url = string
    client_id = string
    client_secret = string
    audience = string
  }))
  default = [
    {
      name = "TF Acceptance Crowdstrike"
      is_active = true
      type = "CrowdStrike"
      sync_frequency = 3600
      endpoint = "https://endpoint"
      token_url = "https://url/token"
      client_id = "client_id"
      client_secret = "client_secret"
      audience = "audience"
    },
    {
      name = "TF Acceptance Wiz"
      is_active = true
      type = "WizIntegrationApi"
      sync_frequency = 3600
      endpoint = "https://endpoint"
      token_url = "https://url/token"
      client_id = "client_id"
      client_secret = "client_secret"
      audience = "audience"
    }
  ]
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
    "crowdstrike" = {
      max_last_seen               = 3600
      match_hostname              = true
      match_serial_number         = true
      prevent_rfm                 = true
    },
    "wiz" = {
      max_last_seen               = 3600
      container_cluster_connected = true
    }
  }
}

resource "aembit_integration" "ac" {
  count = length(var.integrations)
  
  name           = var.integrations[count.index].name
  type           = var.integrations[count.index].type
  sync_frequency = var.integrations[count.index].sync_frequency
  is_active      = var.integrations[count.index].is_active
  endpoint       = var.integrations[count.index].endpoint
  oauth_client_credentials = {
    token_url     = var.integrations[count.index].token_url
    client_id     = var.integrations[count.index].client_id
    client_secret = var.integrations[count.index].client_secret
    audience      = var.integrations[count.index].audience
  }
}

resource "aembit_access_condition" "ac" {
  count = length(var.integrations)

  name           = var.integrations[count.index].name
  is_active      = var.integrations[count.index].is_active
  integration_id = aembit_integration.ac[count.index].id

  dynamic "crowdstrike_conditions" {
    for_each = var.integrations[count.index].type == "CrowdStrike" ? [1] : []
    content {
      max_last_seen               = var.access_conditions["crowdstrike"].max_last_seen
      match_hostname              = var.access_conditions["crowdstrike"].match_hostname
      match_serial_number         = var.access_conditions["crowdstrike"].match_serial_number
      prevent_rfm                 = var.access_conditions["crowdstrike"].prevent_rfm
    }
  }

  dynamic "wiz_conditions" {
    for_each = var.integrations[count.index].type == "WizIntegrationApi" ? [1] : []
    content {
      max_last_seen               = var.access_conditions["wiz"].max_last_seen
      container_cluster_connected = var.access_conditions["wiz"].container_cluster_connected
    }
  }
}