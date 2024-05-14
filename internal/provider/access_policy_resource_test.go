package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_policy.first_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values for First Policy.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "name", "TF First Policy"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "client_workload"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "trust_providers.#", "0"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "access_conditions.#", "0"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "credential_provider"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "server_workload"),

					// Verify values for Second Policy.
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "id"),
					resource.TestCheckResourceAttr("aembit_access_policy.second_policy", "name", "TF Second Policy"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "client_workload"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "credential_provider"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "server_workload"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_policy.first_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values for First Policy.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "name", "Placeholder"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "client_workload"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "trust_providers.#", "0"),
					resource.TestCheckResourceAttr("aembit_access_policy.first_policy", "access_conditions.#", "0"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "credential_provider"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "server_workload"),

					// Verify values for Second Policy.
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "id"),
					resource.TestCheckResourceAttr("aembit_access_policy.second_policy", "name", "Placeholder"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "client_workload"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "credential_provider"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "server_workload"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
