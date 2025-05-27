provider "aembit" {
}

resource "aembit_log_stream" "crowdstrike_http_event_collector" {
	name = "TF Acceptance CrowdstrikeHttpEventCollector LogStream"
	description = "TF Acceptance CrowdstrikeHttpEventCollector LogStream"
	data_type = "AuditLogs"
	type = "CrowdstrikeHttpEventCollector"
	is_active = false
	crowdstrike_http_event_collector = {
		host_port = "crowdstrike.example.com:8088"
		api_key = "c7bbe054cbdf49c3b892de72867fed82"
		source_name = "crowdstrike-test-source"
		tls = false
	}
}
