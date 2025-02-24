provider "aembit" {
}

resource "aembit_access_policy" "multi_cp_duplicate_policy_1" {
		is_active = false
		client_workload = "c460097e-2db7-4190-953d-fddd3a636c71"
		credential_providers = [{
			credential_provider_id = "d939f2f1-8cf2-4296-8f89-81093919f15d",
			mapping_type = "HttpBody",
			httpbody_field_path = "test_field_path_1",
			httpbody_field_value = "test_field_value"
		}, {
			credential_provider_id = "6f88117b-c549-4c3a-867c-55159ae27033",
			mapping_type = "HttpBody",
			httpbody_field_path = "test_field_path_1",
			httpbody_field_value = "test_field_value"
		}]
		server_workload = "eca31347-b739-4522-8628-f78b71e23f8d"
}
