package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testServerWorkloadResource string = "aembit_server_workload.test"

func testDeleteServerWorkload(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteServerWorkload(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccServerWorkloadResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/server/TestAccServerWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/server/TestAccServerWorkloadResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Server Workload Name
					resource.TestCheckResourceAttr(testServerWorkloadResource, "name", "Unit Test 1"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "is_active", "false"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.day", "Sunday"),
					// Verify Service Endpoint.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.host", "unittest.testhost.com"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.port", "443"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.app_protocol", "HTTP"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.transport_protocol", "TCP"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.requested_port", "443"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.tls_verification", "full"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.authentication_config.method", "HTTP Authentication"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.authentication_config.scheme", "Bearer"),
					// Verify HTTP Headers.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.http_headers.%", "3"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.http_headers.host", "graph.microsoft.com"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.http_headers.user-agent", "curl/7.64.1"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.http_headers.accept", "*/*"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testServerWorkloadResource, "id"),
					resource.TestCheckResourceAttrSet(testServerWorkloadResource, "service_endpoint.external_id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testServerWorkloadResource, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: string(createFile), Check: testDeleteServerWorkload(testServerWorkloadResource), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testServerWorkloadResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testServerWorkloadResource, "name", "Unit Test 1 - Modified"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "tags.day", "Tuesday"),
					// Verify Service Endpoint updated.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.host", "unittest.testhost2.com"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.authentication_config.method", "HTTP Authentication"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.authentication_config.scheme", "Header"),
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.authentication_config.config", "X-Vault-Token"),
					// Verify HTTP Headers.
					resource.TestCheckResourceAttr(testServerWorkloadResource, "service_endpoint.http_headers.%", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
