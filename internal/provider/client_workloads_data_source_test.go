package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClientWorkloadsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/data/TestAccClientWorkloadsDataSource.tf")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "unittest1namespace", fmt.Sprintf("unittest1namespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Client Workloads returned
					resource.TestCheckResourceAttrSet("data.aembit_client_workloads.test", "client_workloads.#"),
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.name", "Unit Test 1"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.description", "Acceptance Test client workload"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.identities.#", "1"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.identities.0.type", "k8sNamespace"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.identities.0.value", fmt.Sprintf("unittest1namespace%d", randID)),
					// Verify Tags.
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.tags.color", "blue"),
					resource.TestCheckResourceAttr("data.aembit_client_workloads.test", "client_workloads.0.tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_client_workloads.test", "client_workloads.0.id"),
				),
			},
		},
	})
}
