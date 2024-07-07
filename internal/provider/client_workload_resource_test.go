package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testCWResource string = "aembit_client_workload.test"
const testCWResourceDescription string = "Acceptance Test client workload"
const testCWResourceIdentitiesCount string = "identities.#"
const testCWResourceIdentitiesType string = "identities.0.type"
const testCWResourceIdentitiesValue string = "identities.0.value"

func testDeleteClientWorkload() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[testCWResource]; !ok {
			return fmt.Errorf("Not found: %s", testCWResource)
		}
		if ok, err = testClient.DeleteClientWorkload(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccClientWorkloadResource_k8sNamespace(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/k8sNamespace/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/k8sNamespace/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "unittest1namespace")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1"),
					resource.TestCheckResourceAttr(testCWResource, "description", testCWResourceDescription),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType, "k8sNamespace"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue, newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: createFileConfig, Check: testDeleteClientWorkload(), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_k8sPodName(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/k8sPodName/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/k8sPodName/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "unittest1podname")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - In Resource Set"),
					resource.TestCheckResourceAttr(testCWResource, "description", testCWResourceDescription),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType, "k8sPodName"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue, newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - In Resource Set - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_AwsLambdaArn(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "arn:aws:lambda:us-east-1:880961858887:function:helloworld")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - awsLambdaArn"),
					resource.TestCheckResourceAttr(testCWResource, "description", testCWResourceDescription),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType, "awsLambdaArn"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue, newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - awsLambdaArn - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testCWResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testCWResource, "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
