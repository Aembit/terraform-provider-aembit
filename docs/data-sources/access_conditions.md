---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "aembit_access_conditions Data Source - terraform-provider-aembit"
subcategory: ""
description: |-
  Manages an accessCondition.
---

# aembit_access_conditions (Data Source)

Manages an accessCondition.



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `access_conditions` (Attributes List) List of accessConditions. (see [below for nested schema](#nestedatt--access_conditions))

<a id="nestedatt--access_conditions"></a>
### Nested Schema for `access_conditions`

Read-Only:

- `crowdstrike_conditions` (Attributes) CrowdStrike Specific rules for the Access Condition. (see [below for nested schema](#nestedatt--access_conditions--crowdstrike_conditions))
- `description` (String) User-provided description of the accessCondition.
- `geoip_conditions` (Attributes) Defines geographical conditions for filtering based on country and administrative subdivisions. (see [below for nested schema](#nestedatt--access_conditions--geoip_conditions))
- `id` (String) Unique identifier of the accessCondition.
- `integration_id` (String) ID of the Integration used by the Access Condition.
- `is_active` (Boolean) Active/Inactive status of the accessCondition.
- `name` (String) User-provided name of the accessCondition.
- `tags` (Map of String) Tags are key-value pairs.
- `time_conditions` (Attributes) Defines the conditions for scheduling based on time, including specific time slots and timezone settings for the Access Condition. (see [below for nested schema](#nestedatt--access_conditions--time_conditions))
- `wiz_conditions` (Attributes) Wiz Specific rules for the Access Condition. (see [below for nested schema](#nestedatt--access_conditions--wiz_conditions))

<a id="nestedatt--access_conditions--crowdstrike_conditions"></a>
### Nested Schema for `access_conditions.crowdstrike_conditions`

Required:

- `match_hostname` (Boolean)
- `match_serial_number` (Boolean)
- `max_last_seen` (Number)
- `prevent_rfm` (Boolean)


<a id="nestedatt--access_conditions--geoip_conditions"></a>
### Nested Schema for `access_conditions.geoip_conditions`

Required:

- `locations` (Attributes List) A list of geographical locations, each containing a country code and optional subdivisions. (see [below for nested schema](#nestedatt--access_conditions--geoip_conditions--locations))

<a id="nestedatt--access_conditions--geoip_conditions--locations"></a>
### Nested Schema for `access_conditions.geoip_conditions.locations`

Required:

- `country_code` (String) A list of two-letter country code identifiers (as defined by ISO 3166-1) to allow as part of the validation for this access condition.

Optional:

- `subdivisions` (Attributes List) A list of subdivision identifiers (as defined by ISO 3166) to allow as part of the validation for this access condition. (see [below for nested schema](#nestedatt--access_conditions--geoip_conditions--locations--subdivisions))

<a id="nestedatt--access_conditions--geoip_conditions--locations--subdivisions"></a>
### Nested Schema for `access_conditions.geoip_conditions.locations.subdivisions`

Required:

- `subdivision_code` (String) The subdivision identifier as defined by ISO 3166.




<a id="nestedatt--access_conditions--time_conditions"></a>
### Nested Schema for `access_conditions.time_conditions`

Required:

- `schedule` (Attributes List) (see [below for nested schema](#nestedatt--access_conditions--time_conditions--schedule))
- `timezone` (String) Timezone value such as America/Chicago, Europe/Istanbul

<a id="nestedatt--access_conditions--time_conditions--schedule"></a>
### Nested Schema for `access_conditions.time_conditions.schedule`

Required:

- `day` (String) Day of Week, for example: Tuesday
- `end_time` (String) The end time of the schedule in 24-hour format (HH:mm), e.g., '18:00' for 6:00 PM.
- `start_time` (String) The start time of the schedule in 24-hour format (HH:mm), e.g., '07:00' for 7:00 AM.



<a id="nestedatt--access_conditions--wiz_conditions"></a>
### Nested Schema for `access_conditions.wiz_conditions`

Required:

- `container_cluster_connected` (Boolean)
- `max_last_seen` (Number)
