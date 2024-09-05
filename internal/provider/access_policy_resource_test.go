package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const AccessPolicyPathFirst string = "aembit_access_policy.first_policy"
const AccessPolicyPathSecond string = "aembit_access_policy.second_policy"
const AccessPolicyPathThird string = "aembit_access_policy.third_policy"

var accessPolicyChecks = []resource.TestCheckFunc{
	// Verify values for First Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),
}

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
					accessPolicyChecks...,
				),
			},
			// ImportState testing
			{ResourceName: AccessPolicyPathFirst, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					accessPolicyChecks...,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var basicAccessPolicyChecks = []resource.TestCheckFunc{
	// Verify values for First Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "trust_providers.#", "0"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "access_conditions.#", "0"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "credential_providers.#", "1"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),

	// Verify values for Second Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "client_workload"),
	resource.TestCheckResourceAttr(AccessPolicyPathSecond, "credential_providers.#", "1"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "server_workload"),

	// Third values for Third Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathThird, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathThird, "client_workload"),
	resource.TestCheckResourceAttr(AccessPolicyPathThird, "credential_providers.#", "1"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathThird, "server_workload"),
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "clientworkloadNamespace")
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(createFileConfig, modifyFileConfig, "secondClientWorkloadNamespace")
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(createFileConfig, modifyFileConfig, "thirdClientWorkloadNamespace")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					append(basicAccessPolicyChecks,
						resource.TestCheckResourceAttr(AccessPolicyPathFirst, "name", "TF First Policy"),
						resource.TestCheckResourceAttr(AccessPolicyPathSecond, "name", "TF Second Policy"),
						resource.TestCheckResourceAttr(AccessPolicyPathThird, "name", "TF Third Policy"),
					)...,
				),
			},
			// ImportState testing
			{ResourceName: AccessPolicyPathFirst, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					append(basicAccessPolicyChecks,
						resource.TestCheckResourceAttr(AccessPolicyPathFirst, "name", "Placeholder"),
						resource.TestCheckResourceAttr(AccessPolicyPathSecond, "name", "Placeholder"),
						resource.TestCheckResourceAttr(AccessPolicyPathThird, "name", "Placeholder"),
					)...,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
