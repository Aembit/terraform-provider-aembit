package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testClientWorkloadResource string = "aembit_client_workload.test"

func testDeleteClientWorkload() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[testClientWorkloadResource]; !ok {
			return fmt.Errorf("Not found: %s", testClientWorkloadResource)
		}

		if ok, err = Client.DeleteClientWorkload(rs.Primary.ID, nil); !ok {
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
					resource.TestCheckResourceAttr(testClientWorkloadResource, "name", "Unit Test 1"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "description", "Acceptance Test client workload"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.#", "1"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.0.type", "k8sNamespace"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.0.value", newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testClientWorkloadResource, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             createFileConfig,
				Check:              testDeleteClientWorkload(),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{
				Config: createFileConfig,
			},
			// ImportState testing
			{
				ResourceName:      testClientWorkloadResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testClientWorkloadResource, "name", "Unit Test 1 - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_AwsLambdaArn(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testClientWorkloadResource, "name", "Unit Test 1 - awsLambdaArn"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "description", "Acceptance Test client workload"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.#", "1"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.0.type", "awsLambdaArn"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "identities.0.value", "arn:aws:lambda:us-east-1:880961858887:function:helloworld"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.color", "blue"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testClientWorkloadResource, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testClientWorkloadResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testClientWorkloadResource, "name", "Unit Test 1 - awsLambdaArn - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.%", "2"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.color", "orange"),
					resource.TestCheckResourceAttr(testClientWorkloadResource, "tags.day", "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
