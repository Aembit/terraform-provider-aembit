package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testCWResource string = "aembit_client_workload.test"
const testCWResourceDescription string = "Acceptance Test client workload"
const testCWResourceIdentitiesCount string = "identities.#"

var testCWResourceIdentitiesType = []string{"identities.0.type", "identities.1.type", "identities.2.type", "identities.3.type"}
var testCWResourceIdentitiesValue = []string{"identities.0.value", "identities.1.value", "identities.2.value", "identities.3.value"}

func testDeleteClientWorkload() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[testCWResource]; !ok {
			return fmt.Errorf("Not found: %s", testCWResource)
		}
		if ok, err = testClient.DeleteClientWorkload(context.Background(), rs.Primary.ID, nil); !ok {
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
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "k8sNamespace"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
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
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_k8sPodName(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/k8sPodName/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/k8sPodName/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newNamePod1 := randomizeFileConfigs(string(createFile), string(modifyFile), "unittest1podname1")
	createFileConfig, modifyFileConfig, newNamePod2 := randomizeFileConfigs(createFileConfig, modifyFileConfig, "unittest1podname2")

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
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "k8sPodName"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], newNamePod1),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[1], "k8sPodName"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[1], newNamePod2),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
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
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccClientWorkloadResource_k8sPodName_CustomResourceSet tests resource creation within a custom resource set.
func TestAccClientWorkloadResource_k8sPodName_CustomResourceSetAuth(t *testing.T) {
	skipNotCI(t)

	createFile, _ := os.ReadFile("../../tests/client/resourceSet/TestAccClientWorkloadCustomResourceSet.tf")
	createFileConfig, _, newName := randomizeFileConfigs(string(createFile), "", "custom-resource-set")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "TF Acceptance RS"),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "k8sPodName"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
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
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "awsLambdaArn"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], newName),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
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
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_GitLabJob(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/gitLabJob/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/gitLabJob/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newSubject := randomizeFileConfigs(string(createFile), string(modifyFile), "subject")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - gitLabJob"),
					resource.TestCheckResourceAttr(testCWResource, "description", testCWResourceDescription),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "4"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "gitlabIdTokenNamespacePath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], "namespacePath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[1], "gitlabIdTokenProjectPath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[1], "projectPath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[2], "gitlabIdTokenRefPath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[2], "refPath"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[3], "gitlabIdTokenSubject"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[3], newSubject),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
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
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - gitLabJob - modified"),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_StandaloneCA(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/client/standalone-certificate-authority/TestAccClientWorkloadStandaloneCertificateAuthority.tf")
	createFileConfig, _, newName := randomizeFileConfigs(string(createFile), "", "unittestname")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", newName),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], "k8sPodName"),
					resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], newName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
					resource.TestCheckResourceAttrSet(testCWResource, "standalone_certificate_authority"),
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var clientWorkloadIdentifierTests = []struct {
	identityType  string
	identityValue string
}{
	{"awsAccountId", "123456789012"},
	{"awsRegion", "us-east-1"},
	{"awsEc2InstanceId", "i-0b22a22eec53b9321"},
	{"azureSubscriptionId", "e58ac327-32b9-414f-a3f5-50850d4dc343"},
	{"azureVmId", "4ce4ead0-0561-4d8c-8313-16ebeb11c1b2"},
}

func TestAccClientWorkloadResource_Miscellaneous(t *testing.T) {
	createFileConfigWithPlaceholders, _ := os.ReadFile("../../tests/client/miscellaneous/TestAccClientWorkloadResource.tf")
	modifyFileConfigWithPlaceholders, _ := os.ReadFile("../../tests/client/miscellaneous/TestAccClientWorkloadResource.tfmod")

	for _, test := range clientWorkloadIdentifierTests {

		createFileConfig := strings.ReplaceAll(string(createFileConfigWithPlaceholders), "IDENTITY_TYPE_PLACEHOLDER", test.identityType)
		createFileConfig = strings.ReplaceAll(createFileConfig, "IDENTITY_VALUE_PLACEHOLDER", test.identityValue)

		modifyFileConfig := strings.ReplaceAll(string(modifyFileConfigWithPlaceholders), "IDENTITY_TYPE_PLACEHOLDER", test.identityType)
		modifyFileConfig = strings.ReplaceAll(modifyFileConfig, "IDENTITY_VALUE_PLACEHOLDER", test.identityValue)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: createFileConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Client Workload Name, Description, Active status
						resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - miscellaneous"),
						resource.TestCheckResourceAttr(testCWResource, "description", testCWResourceDescription),
						resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
						// Verify Workload Identity.
						resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesCount, "1"),
						resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesType[0], test.identityType),
						resource.TestCheckResourceAttr(testCWResource, testCWResourceIdentitiesValue[0], test.identityValue),
						// Verify Tags.
						resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
						resource.TestCheckResourceAttr(testCWResource, tagsColor, "blue"),
						resource.TestCheckResourceAttr(testCWResource, tagsDay, "Sunday"),
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
						resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1 - miscellaneous - modified"),
						// Verify active state updated.
						resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
						// Verify Tags.
						resource.TestCheckResourceAttr(testCWResource, tagsCount, "2"),
						resource.TestCheckResourceAttr(testCWResource, tagsColor, "orange"),
						resource.TestCheckResourceAttr(testCWResource, tagsDay, "Tuesday"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}
