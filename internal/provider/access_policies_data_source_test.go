package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessPoliciesDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/data/TestAccAccessPoliciesDataSource.tf")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Access Policies returned
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.#"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.id"),
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.client_workload"),
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.trust_providers.#"),
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.access_conditions.#"),
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.credential_provider"),
					resource.TestCheckResourceAttrSet("data.aembit_access_policies.test", "access_policies.0.server_workload"),
				),
			},
		},
	})
}
