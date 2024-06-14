package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerWorkloadsDataSource(t *testing.T) {

	createFile, _ := os.ReadFile("../../tests/server/data/TestAccServerWorkloadsDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Server Workloads returned
					resource.TestCheckResourceAttrSet("data.aembit_server_workloads.test", "server_workloads.#"),
					// Verify the attributes of the first Server Workload
					// Verify Server Workload Name
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.name", "Unit Test 1"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.tags.color", "blue"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.tags.day", "Sunday"),
					// Verify Service Endpoint.
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.host", "unittest.testhost.com"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.port", "443"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.app_protocol", "HTTP"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.transport_protocol", "TCP"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.requested_port", "443"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.tls_verification", "full"),
					// Service Endpoint Authentication config is not returned by the API
					//resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.authentication_config.method", "HTTP Authentication"),
					//resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.authentication_config.scheme", "Bearer"),
					// Verify HTTP Headers.
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.http_headers.%", "3"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.http_headers.host", "graph.microsoft.com"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.http_headers.user-agent", "curl/7.64.1"),
					resource.TestCheckResourceAttr("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.http_headers.accept", "*/*"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_server_workloads.test", "server_workloads.0.id"),
					resource.TestCheckResourceAttrSet("data.aembit_server_workloads.test", "server_workloads.0.service_endpoint.external_id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("data.aembit_server_workloads.test", "server_workloads.0.id"),
				),
			},
		},
	})
}
