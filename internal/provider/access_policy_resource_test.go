package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const AccessPolicyPathFirst string = "aembit_access_policy.first_policy"
const AccessPolicyPathSecond string = "aembit_access_policy.second_policy"
const AccessPolicyPathMultiCPFirst string = "aembit_access_policy.multi_cp_first_policy"
const AccessPolicyPathMultiCPSecond string = "aembit_access_policy.multi_cp_second_policy"

var accessPolicyChecks = []resource.TestCheckFunc{
	// Verify values for First Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "credential_provider"),
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
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),

	// Verify values for Second Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "client_workload"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "credential_provider"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "server_workload"),
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "clientworkloadNamespace")
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(createFileConfig, modifyFileConfig, "secondClientWorkloadNamespace")

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
					)...,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMultipleCredentialProviders_AccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(string(createFile), string(modifyFile), "clientworkloadNamespace")
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(createFileConfig, modifyFileConfig, "secondClientWorkloadNamespace")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPFirst, "name", "TF Multi CP First Policy"),
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPFirst, "credential_providers.#", "2"),

					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPSecond, "name", "TF Multi CP Second Policy"),
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPSecond, "credential_providers.#", "3"),
				),
			},
			// ImportState testing
			{ResourceName: AccessPolicyPathMultiCPFirst, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPFirst, "name", "TF Multi CP First Policy Updated"),
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPFirst, "credential_providers.#", "2"),

					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPSecond, "name", "TF Multi CP Second Policy Updated"),
					resource.TestCheckResourceAttr(AccessPolicyPathMultiCPSecond, "credential_providers.#", "2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
