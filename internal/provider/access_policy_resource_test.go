package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const AccessPolicyPathFirst string = "aembit_access_policy.first_policy"
const AccessPolicyPathSecond string = "aembit_access_policy.second_policy"

func TestAccAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "clientworkloadNamespace")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      AccessPolicyPathFirst,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "clientworkloadNamespace")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values for First Policy.
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "name", "TF First Policy"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "trust_providers.#", "0"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "access_conditions.#", "0"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "credential_provider"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),

					// Verify values for Second Policy.
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "id"),
					resource.TestCheckResourceAttr(AccessPolicyPathSecond, "name", "TF Second Policy"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "client_workload"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "credential_provider"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "server_workload"),
				),
			},
			// ImportState testing
			{
				ResourceName:      AccessPolicyPathFirst,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values for First Policy.
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "name", "Placeholder"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "trust_providers.#", "0"),
					resource.TestCheckResourceAttr(AccessPolicyPathFirst, "access_conditions.#", "0"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "credential_provider"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),

					// Verify values for Second Policy.
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "id"),
					resource.TestCheckResourceAttr(AccessPolicyPathSecond, "name", "Placeholder"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "client_workload"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "credential_provider"),
					resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "server_workload"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
