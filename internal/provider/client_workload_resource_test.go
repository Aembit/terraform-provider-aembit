package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCWResource                string = "aembit_client_workload.test"
	testCWResourceDescription     string = "Acceptance Test client workload"
	testCWResourceIdentitiesCount string = "identities.#"
)

var (
	testCWResourceIdentitiesType = []string{
		"identities.0.type",
		"identities.1.type",
		"identities.2.type",
		"identities.3.type",
	}
	testCWResourceIdentitiesValue = []string{
		"identities.0.value",
		"identities.1.value",
		"identities.2.value",
		"identities.3.value",
	}

	testCWResourceIdentitiesClaimName = []string{
		"identities.0.claim_name",
		"identities.1.claim_name",
		"identities.2.claim_name",
		"identities.3.claim_name",
	}

	testCWResourceIdentitiesIdpId = []string{
		"identities.0.identity_provider_id",
		"identities.1.identity_provider_id",
		"identities.2.identity_provider_id",
		"identities.3.identity_provider_id",
	}
)

func testDeleteClientWorkload() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[testCWResource]; !ok {
			return fmt.Errorf("Not found: %s", testCWResource)
		}
		resourceSetId := getResourceSetId(types.StringValue(rs.Primary.Attributes["resource_set_id"]), testClient)
		if ok, err = testClient.DeleteClientWorkload(context.Background(), rs.Primary.ID, nil, &resourceSetId); !ok {
			return err
		}
		return nil
	}
}

func TestAccClientWorkloadResource_k8sNamespace(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/client/k8sNamespace/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/k8sNamespace/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"unittest1namespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(testCWResource, "name", "Unit Test 1"),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						testCWResourceDescription,
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"k8sNamespace",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_k8sPodName(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/client/k8sPodName/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/k8sPodName/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newNamePod1 := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"unittest1podname1",
	)
	createFileConfig, modifyFileConfig, newNamePod2 := randomizeFileConfigs(
		createFileConfig,
		modifyFileConfig,
		"unittest1podname2",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - In Resource Set",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						testCWResourceDescription,
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"k8sPodName",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newNamePod1,
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[1],
						"k8sPodName",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[1],
						newNamePod2,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - In Resource Set - modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_AwsLambdaArn(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/awsLambdaArn/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"arn:aws:lambda:us-east-1:880961858887:function:helloworld",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - awsLambdaArn",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						testCWResourceDescription,
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"awsLambdaArn",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - awsLambdaArn - modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_GitLabJob(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/client/gitLabJob/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/client/gitLabJob/TestAccClientWorkloadResource.tfmod")
	createFileConfig, modifyFileConfig, newSubject := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"subject",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - gitLabJob",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						testCWResourceDescription,
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"4",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"gitlabIdTokenNamespacePath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						"namespacePath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[1],
						"gitlabIdTokenProjectPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[1],
						"projectPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[2],
						"gitlabIdTokenRefPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[2],
						"refPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[3],
						"gitlabIdTokenSubject",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[3],
						newSubject,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"Unit Test 1 - gitLabJob - modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_StandaloneCA(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/standalone-certificate-authority/TestAccClientWorkloadStandaloneCertificateAuthority.tf",
	)
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
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"k8sPodName",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
					resource.TestCheckResourceAttrSet(
						testCWResource,
						"standalone_certificate_authority",
					),
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
	t.Parallel()
	createFileConfigWithPlaceholders, _ := os.ReadFile(
		"../../tests/client/miscellaneous/TestAccClientWorkloadResource.tf",
	)
	modifyFileConfigWithPlaceholders, _ := os.ReadFile(
		"../../tests/client/miscellaneous/TestAccClientWorkloadResource.tfmod",
	)

	for _, test := range clientWorkloadIdentifierTests {
		createFileConfig := strings.ReplaceAll(
			string(createFileConfigWithPlaceholders),
			"IDENTITY_TYPE_PLACEHOLDER",
			test.identityType,
		)
		createFileConfig = strings.ReplaceAll(
			createFileConfig,
			"IDENTITY_VALUE_PLACEHOLDER",
			test.identityValue,
		)

		modifyFileConfig := strings.ReplaceAll(
			string(modifyFileConfigWithPlaceholders),
			"IDENTITY_TYPE_PLACEHOLDER",
			test.identityType,
		)
		modifyFileConfig = strings.ReplaceAll(
			modifyFileConfig,
			"IDENTITY_VALUE_PLACEHOLDER",
			test.identityValue,
		)

		createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(
			createFileConfig,
			modifyFileConfig,
			"unittest1podname1",
		)

		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Create and Read testing
				{
					Config: createFileConfig,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify Client Workload Name, Description, Active status
						resource.TestCheckResourceAttr(
							testCWResource,
							"name",
							"Unit Test 1 - miscellaneous",
						),
						resource.TestCheckResourceAttr(
							testCWResource,
							"description",
							testCWResourceDescription,
						),
						resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
						// Verify Workload Identity.
						resource.TestCheckResourceAttr(
							testCWResource,
							testCWResourceIdentitiesCount,
							"2",
						),
						resource.TestCheckResourceAttr(
							testCWResource,
							testCWResourceIdentitiesType[0],
							test.identityType,
						),
						resource.TestCheckResourceAttr(
							testCWResource,
							testCWResourceIdentitiesValue[0],
							test.identityValue,
						),
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
						resource.TestCheckResourceAttr(
							testCWResource,
							"name",
							"Unit Test 1 - miscellaneous - modified",
						),
						// Verify active state updated.
						resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					),
				},
				// Delete testing automatically occurs in TestCase
			},
		})
	}
}

func TestAccClientWorkloadResource_ProcessPath(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/processPath/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/processPath/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"/process/path",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance ProcessPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"processPath",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance ProcessPath - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_ProcessCommandLine(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/processCommandLine/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/processCommandLine/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"*process command line*",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance ProcessCommandLine",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"processCommandLine",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance ProcessCommandLine - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_OauthRedirectUri(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/oauthRedirectUri/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/oauthRedirectUri/TestAccClientWorkloadResource.tfmod",
	)

	createFileConfig, modifyFileConfig, newValue := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"replace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oauth Redirect URI",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"oauthRedirectUri",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						fmt.Sprintf("https://%s.aembit.local:12345", newValue),
					),
					resource.TestCheckResourceAttr(testCWResource, "enforce_sso", "true"),
					resource.TestCheckResourceAttr(
						testCWResource,
						"sso_identity_providers.#",
						"2",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
					resource.TestCheckResourceAttrSet(
						testCWResource,
						"sso_identity_providers.0",
					),
					resource.TestCheckResourceAttrSet(
						testCWResource,
						"sso_identity_providers.1",
					),
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oauth Redirect URI - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
					resource.TestCheckResourceAttr(testCWResource, "enforce_sso", "true"),
					resource.TestCheckResourceAttr(
						testCWResource,
						"sso_identity_providers.#",
						"1",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_MixedRedirectURIIdentitiesReturnsPlanError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_client_workload" "test" {
  name        = "TF Acceptance - Invalid Mixed Redirect URI"
  description = "Acceptance Test Client Workload"
  is_active   = false
  identities = [
    {
      type  = "oauthRedirectUri"
      value = "https://example.com/callback"
    },
    {
      type  = "k8sNamespace"
      value = "default"
    }
  ]
}
`,
				ExpectError: regexp.MustCompile("oauthRedirectUri"),
			},
		},
	})
}

func TestAccClientWorkloadResource_RedirectURIWithSSOAndMissingIDPsReturnsPlanError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_client_workload" "test" {
  name        = "TF Acceptance - Missing Redirect URI SSO IDPs"
  description = "Acceptance Test Client Workload"
  is_active   = false
  enforce_sso = true
  identities = [
    {
      type  = "oauthRedirectUri"
      value = "https://example.com/callback"
    }
  ]
}
`,
				ExpectError: regexp.MustCompile("sso_identity_providers"),
			},
		},
	})
}

func TestAccClientWorkloadResource_RedirectURIWithDisabledSSOAndIDPsReturnsPlanError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_client_workload" "test" {
  name        = "TF Acceptance - Redirect URI Disabled SSO With IDPs"
  description = "Acceptance Test Client Workload"
  is_active   = false
  enforce_sso = false
  identities = [
    {
      type  = "oauthRedirectUri"
      value = "https://example.com/callback"
    }
  ]
  sso_identity_providers = [
    "11111111-1111-1111-1111-111111111111",
  ]
}
`,
				ExpectError: regexp.MustCompile("sso_identity_providers"),
			},
		},
	})
}

func TestAccClientWorkloadResource_NonRedirectWithFalseEnforceSSOReturnsPlanError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_client_workload" "test" {
  name        = "TF Acceptance - Invalid Enforce SSO"
  description = "Acceptance Test Client Workload"
  is_active   = false
  enforce_sso = false
  identities = [
    {
      type  = "k8sNamespace"
      value = "tf-invalid-enforce-sso"
    }
  ]
}
`,
				ExpectError: regexp.MustCompile("enforce_sso"),
			},
		},
	})
}

func TestAccClientWorkloadResource_NonRedirectWithIDPsReturnsPlanError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "aembit_client_workload" "test" {
  name        = "TF Acceptance - Invalid SSO IDP Mapping"
  description = "Acceptance Test Client Workload"
  is_active   = false
  identities = [
    {
      type  = "k8sNamespace"
      value = "tf-invalid-sso-idp-mapping"
    }
  ]
  sso_identity_providers = [
    "11111111-1111-1111-1111-111111111111",
  ]
}
`,
				ExpectError: regexp.MustCompile("sso_identity_providers"),
			},
		},
	})
}

// TestConvertClientWorkloadDTOToModel_NormalizesSSOFields verifies that SSO fields are normalized correctly when converting DTO to Model.
func TestConvertClientWorkloadDTOToModel_NormalizesSSOFields(t *testing.T) {
	t.Parallel()

	redirectModel := convertClientWorkloadDTOToModel(context.Background(), aembit.ClientWorkloadExternalDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "4fce7a89-9ab2-4fbf-a757-acbc8a1338d4",
				Name:       "redirect workload",
				IsActive:   true,
			},
		},
		Identities: []aembit.ClientWorkloadIdentityDTO{
			{Type: "oauthRedirectUri", Value: "https://example.com/callback"},
		},
		EnforceSso:           false,
		SsoIdentityProviders: []string{"11111111-1111-1111-1111-111111111111"},
	}, &models.ClientWorkloadResourceModel{})

	assertSetContainsString(t, redirectModel.SsoIdentityProviders, "11111111-1111-1111-1111-111111111111")
	if redirectModel.EnforceSso.ValueBool() {
		t.Fatalf("expected redirect workload enforce_sso to stay false")
	}

	nonRedirectModel := convertClientWorkloadDTOToModel(context.Background(), aembit.ClientWorkloadExternalDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "5fce7a89-9ab2-4fbf-a757-acbc8a1338d4",
				Name:       "non-redirect workload",
				IsActive:   true,
			},
		},
		Identities: []aembit.ClientWorkloadIdentityDTO{
			{Type: "k8sNamespace", Value: "default"},
		},
		EnforceSso: false,
	}, &models.ClientWorkloadResourceModel{})

	if !nonRedirectModel.EnforceSso.ValueBool() {
		t.Fatalf("expected non-redirect workload enforce_sso to normalize to true")
	}
	if !nonRedirectModel.SsoIdentityProviders.IsNull() {
		t.Fatalf("expected non-redirect workload sso_identity_providers to be null")
	}
}

// TestConvertClientWorkloadDTOToModel_NormalizesCimdSSOFields verifies that cimdClientId preserves enforce_sso and sso_identity_providers from DTO.
func TestConvertClientWorkloadDTOToModel_NormalizesCimdSSOFields(t *testing.T) {
	t.Parallel()

	cimdModel := convertClientWorkloadDTOToModel(context.Background(), aembit.ClientWorkloadExternalDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "4fce7a89-9ab2-4fbf-a757-acbc8a1338d4",
				Name:       "cimd workload",
				IsActive:   true,
			},
		},
		Identities: []aembit.ClientWorkloadIdentityDTO{
			{Type: "cimdClientId", Value: "https://mcpjam.com/.well-known/oauth/client-metadata.json"},
		},
		EnforceSso:           false,
		SsoIdentityProviders: []string{"11111111-1111-1111-1111-111111111111"},
	}, &models.ClientWorkloadResourceModel{})

	assertSetContainsString(t, cimdModel.SsoIdentityProviders, "11111111-1111-1111-1111-111111111111")
	if cimdModel.EnforceSso.ValueBool() {
		t.Fatalf("expected cimd workload enforce_sso to stay false")
	}
}

// TestConvertClientWorkloadModelToDTO_ValidCimdWithDisabledSSO verifies valid cimdClientId with disabled enforce_sso.
func TestConvertClientWorkloadModelToDTO_ValidCimdWithDisabledSSO(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("cimdClientId"), Value: types.StringValue("https://mcpjam.com/.well-known/oauth/client-metadata.json"), ClaimName: types.StringNull()}},
		types.BoolValue(false),
		types.SetNull(types.StringType),
	)

	dto, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if dto.EnforceSso {
		t.Fatalf("expected enforce_sso to remain false for cimd client ID workload")
	}
}

// TestConvertClientWorkloadModelToDTO_ValidCimdWithSSOIdentityProviders verifies valid cimdClientId with enforce_sso and sso identity providers.
func TestConvertClientWorkloadModelToDTO_ValidCimdWithSSOIdentityProviders(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("cimdClientId"), Value: types.StringValue("https://mcpjam.com/.well-known/oauth/client-metadata.json"), ClaimName: types.StringNull()}},
		types.BoolValue(true),
		types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")}),
	)

	dto, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(dto.SsoIdentityProviders) != 1 {
		t.Fatalf("expected one sso_identity_provider, got %v", dto.SsoIdentityProviders)
	}
}

// TestConvertClientWorkloadModelToDTO_CimdWithSSOAndMissingIDPsReturnsError verifies error when cimdClientId has enforce_sso true but missing IDPs.
func TestConvertClientWorkloadModelToDTO_CimdWithSSOAndMissingIDPsReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("cimdClientId"), Value: types.StringValue("https://mcpjam.com/.well-known/oauth/client-metadata.json"), ClaimName: types.StringNull()}},
		types.BoolValue(true),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when enforce_sso is true and no idps are set for cimdClientId")
	}
}

// TestConvertClientWorkloadModelToDTO_CimdWithDisabledSSOAndIdentityProvidersReturnsError verifies error when cimdClientId has disabled SSO and IDPs configured.
func TestConvertClientWorkloadModelToDTO_CimdWithDisabledSSOAndIdentityProvidersReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("cimdClientId"), Value: types.StringValue("https://mcpjam.com/.well-known/oauth/client-metadata.json"), ClaimName: types.StringNull()}},
		types.BoolValue(false),
		types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")}),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when enforce_sso is false and idps are set for cimdClientId")
	}
}

// TestConvertClientWorkloadModelToDTO_MixedCimdIdentitiesReturnsDiagnostics verifies error when cimdClientId is mixed with other non-MCP identity types.
func TestConvertClientWorkloadModelToDTO_MixedCimdIdentitiesReturnsDiagnostics(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{
			{Type: types.StringValue("cimdClientId"), Value: types.StringValue("https://mcpjam.com/.well-known/oauth/client-metadata.json"), ClaimName: types.StringNull()},
			{Type: types.StringValue("k8sNamespace"), Value: types.StringValue("default"), ClaimName: types.StringNull()},
		},
		types.BoolValue(true),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for mixed cimd client ID identities")
	}
}

// TestConvertClientWorkloadModelToDTO_InvalidCimdUrlReturnsDiagnostics verifies error when cimdClientId has an invalid URL.
func TestConvertClientWorkloadModelToDTO_InvalidCimdUrlReturnsDiagnostics(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{
			{Type: types.StringValue("cimdClientId"), Value: types.StringValue("http://not-https.com/path"), ClaimName: types.StringNull()},
			{Type: types.StringNull(), Value: types.StringNull(), ClaimName: types.StringNull()},
		},
		types.BoolValue(false),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for invalid CIMD Client ID URL")
	}
}

// disabled enforce_sso should not have IDPs.
func TestConvertClientWorkloadModelToDTO_ValidRedirectURIWithDisabledSSO(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("oauthRedirectUri"), Value: types.StringValue("https://example.com/callback"), ClaimName: types.StringNull()}},
		types.BoolValue(false),
		types.SetNull(types.StringType),
	)

	dto, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if dto.EnforceSso {
		t.Fatalf("expected enforce_sso to remain false for redirect URI workload")
	}
}

func TestConvertClientWorkloadModelToDTO_ValidRedirectURIWithSSOIdentityProviders(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("oauthRedirectUri"), Value: types.StringValue("https://example.com/callback"), ClaimName: types.StringNull()}},
		types.BoolValue(true),
		types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")}),
	)

	dto, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(dto.SsoIdentityProviders) != 1 {
		t.Fatalf("expected one sso_identity_provider, got %v", dto.SsoIdentityProviders)
	}
}

// if oauthRedirectUri is specified, then no other identity types are allowed.
func TestConvertClientWorkloadModelToDTO_MixedRedirectURIIdentitiesReturnsDiagnostics(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{
			{Type: types.StringValue("oauthRedirectUri"), Value: types.StringValue("https://example.com/callback"), ClaimName: types.StringNull()},
			{Type: types.StringValue("k8sNamespace"), Value: types.StringValue("default"), ClaimName: types.StringNull()},
		},
		types.BoolValue(true),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics for mixed redirect URI identities")
	}
}

func TestConvertClientWorkloadModelToDTO_RedirectURIWithSSOAndMissingIDPsReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("oauthRedirectUri"), Value: types.StringValue("https://example.com/callback"), ClaimName: types.StringNull()}},
		types.BoolValue(true),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when enforce_sso is true and no idps are set")
	}
}

func TestConvertClientWorkloadModelToDTO_RedirectURIWithDisabledSSOAndIdentityProvidersReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("oauthRedirectUri"), Value: types.StringValue("https://example.com/callback"), ClaimName: types.StringNull()}},
		types.BoolValue(false),
		types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")}),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when enforce_sso is false and idps are set")
	}
}

func TestConvertClientWorkloadModelToDTO_NonRedirectURIWithFalseEnforceSSOReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("k8sNamespace"), Value: types.StringValue("default"), ClaimName: types.StringNull()}},
		types.BoolValue(false),
		types.SetNull(types.StringType),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when non-redirect workload sets enforce_sso to false")
	}
}

func TestConvertClientWorkloadModelToDTO_NonRedirectURIWithIdentityProvidersReturnsError(t *testing.T) {
	t.Parallel()

	model := newClientWorkloadResourceModel(
		t,
		[]models.IdentitiesModel{{Type: types.StringValue("k8sNamespace"), Value: types.StringValue("default"), ClaimName: types.StringNull()}},
		types.BoolValue(true),
		types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")}),
	)

	_, diags := convertClientWorkloadModelToDTO(context.Background(), model, nil, nil)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics when non-redirect workload sets sso_identity_providers")
	}
}

func assertSetContainsString(t *testing.T, set types.Set, expected string) {
	t.Helper()

	var values []string
	diags := set.ElementsAs(context.Background(), &values, false)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics reading set: %v", diags)
	}

	for _, value := range values {
		if value == expected {
			return
		}
	}

	t.Fatalf("expected set to contain %q, got %v", expected, values)
}

func newClientWorkloadResourceModel(
	t *testing.T,
	identities []models.IdentitiesModel,
	enforceSso types.Bool,
	ssoIdentityProviders types.Set,
) models.ClientWorkloadResourceModel {
	t.Helper()

	identitiesSet := types.SetNull(models.TfIdentityObjectType)
	if identities != nil {
		identitiesSet = types.SetValueMust(models.TfIdentityObjectType, identityModelsToAttrValues(identities))
	}

	return models.ClientWorkloadResourceModel{
		Name:                           types.StringValue("test workload"),
		Description:                    types.StringValue("test description"),
		IsActive:                       types.BoolValue(true),
		Identities:                     identitiesSet,
		EnforceSso:                     enforceSso,
		SsoIdentityProviders:           ssoIdentityProviders,
		Tags:                           types.MapNull(types.StringType),
		StandaloneCertificateAuthority: types.StringNull(),
	}
}

func identityModelsToAttrValues(identities []models.IdentitiesModel) []attr.Value {
	values := make([]attr.Value, len(identities))
	for i, identity := range identities {
		idpId := identity.IdentityProviderId
		if idpId.IsNull() || idpId.IsUnknown() {
			idpId = types.StringNull()
		}
		claimName := identity.ClaimName
		if claimName.IsNull() || claimName.IsUnknown() {
			claimName = types.StringNull()
		}
		values[i] = types.ObjectValueMust(models.TfIdentityObjectType.AttrTypes, map[string]attr.Value{
			"type":                 identity.Type,
			"value":                identity.Value,
			"claim_name":           claimName,
			"identity_provider_id": idpId,
		})
	}

	return values
}

func TestAccClientWorkloadResource_OauthScope(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/oauthScope/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/oauthScope/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"*oauth scope*",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oauth Scope",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"oauthScope",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oauth Scope - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccClientWorkloadResource_OidcIdToken(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/oidcIdToken/TestAccClientWorkloadOidcIdToken.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/oidcIdToken/TestAccClientWorkloadOidcIdToken.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"*claim_value*",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oidc Id Token",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"oidcIdToken",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesClaimName[0],
						"claim_key",
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - Oidc Id Token - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccClientWorkloadResource_CimdClientId verifies acceptance testing of a client workload configured with cimdClientId identity.
func TestAccClientWorkloadResource_CimdClientId(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "aembit_client_workload" "test" {
						name        = "TF Acceptance - CIMD Client ID"
						description = "Acceptance Test client workload with CIMD URL"
						enforce_sso = false
						identities = [
							{
								type  = "cimdClientId"
								value = "https://mcpjam.com/.well-known/oauth/client-metadata.json"
							}
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testCWResource, "name", "TF Acceptance - CIMD Client ID"),
					resource.TestCheckResourceAttr(testCWResource, "description", "Acceptance Test client workload with CIMD URL"),
					resource.TestCheckResourceAttr(testCWResource, "enforce_sso", "false"),
					resource.TestCheckResourceAttr(testCWResource, "identities.0.type", "cimdClientId"),
					resource.TestCheckResourceAttr(testCWResource, "identities.0.value", "https://mcpjam.com/.well-known/oauth/client-metadata.json"),
					resource.TestCheckResourceAttrSet(testCWResource, "id"),
				),
			},
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestAccClientWorkloadResource_CimdClientId_Invalid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "aembit_client_workload" "test" {
						name        = "TF Acceptance - CIMD Client ID Invalid"
						description = "Acceptance Test client workload with invalid CIMD URL"
						identities = [
							{
								type  = "cimdClientId"
								value = "http://mcpjam.com" // Needs https:// and a path
							}
						]
					}
				`,
				ExpectError: regexp.MustCompile("Invalid CIMD Client ID Value"),
			},
		},
	})
}

// TestAccClientWorkloadResource_UserGroup verifies acceptance testing of a client workload configured with userGroup identity.
func TestAccClientWorkloadResource_UserGroup(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/client/userGroup/TestAccClientWorkloadResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/client/userGroup/TestAccClientWorkloadResource.tfmod",
	)
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"*group_name*",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Client Workload Name, Description, Active status
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - UserGroup",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						"description",
						"Acceptance Test Client Workload",
					),
					resource.TestCheckResourceAttr(testCWResource, "is_active", "false"),
					// Verify Workload Identity.
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesCount,
						"1",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesType[0],
						"userGroup",
					),
					resource.TestCheckResourceAttr(
						testCWResource,
						testCWResourceIdentitiesValue[0],
						newName,
					),
					resource.TestCheckResourceAttrSet(
						testCWResource,
						testCWResourceIdentitiesIdpId[0],
					),
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
					resource.TestCheckResourceAttr(
						testCWResource,
						"name",
						"TF Acceptance - UserGroup - Modified",
					),
					// Verify active state updated.
					resource.TestCheckResourceAttr(testCWResource, "is_active", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestConvertClientWorkloadModelToDTO_UserGroupIdentity verifies identity_provider_id maps to DTO Key.
func TestConvertClientWorkloadModelToDTO_UserGroupIdentity(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	idpID := "11111111-2222-3333-4444-555555555555"
	groupVal := "engineering"

	identities := []models.IdentitiesModel{
		{
			Type:               types.StringValue("userGroup"),
			Value:              types.StringValue(groupVal),
			IdentityProviderId: types.StringValue(idpID),
			ClaimName:          types.StringNull(),
		},
	}
	identitiesSet, diags := types.SetValueFrom(ctx, models.TfIdentityObjectType, identities)
	if diags.HasError() {
		t.Fatalf("unexpected error creating identities set: %v", diags)
	}

	model := models.ClientWorkloadResourceModel{
		Name:        types.StringValue("Test CW UserGroup"),
		Description: types.StringValue("Test description"),
		IsActive:    types.BoolValue(true),
		Identities:  identitiesSet,
		EnforceSso:  types.BoolValue(true),
		Tags:        types.MapNull(types.StringType),
	}

	dto, diags := convertClientWorkloadModelToDTO(ctx, model, nil, nil)
	if diags.HasError() {
		t.Fatalf("unexpected error converting model to DTO: %v", diags)
	}

	if len(dto.Identities) != 1 {
		t.Fatalf("expected 1 identity in DTO, got %d", len(dto.Identities))
	}

	if dto.Identities[0].Type != "userGroup" {
		t.Errorf("expected identity Type 'userGroup', got '%s'", dto.Identities[0].Type)
	}
	if dto.Identities[0].Value != groupVal {
		t.Errorf("expected identity Value '%s', got '%s'", groupVal, dto.Identities[0].Value)
	}
	if dto.Identities[0].Key != idpID {
		t.Errorf("expected identity Key '%s', got '%s'", idpID, dto.Identities[0].Key)
	}
}

// TestConvertClientWorkloadDTOToModel_UserGroupIdentity verifies DTO Key maps to IdentityProviderId.
func TestConvertClientWorkloadDTOToModel_UserGroupIdentity(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	idpID := "11111111-2222-3333-4444-555555555555"
	groupVal := "engineering"

	dto := aembit.ClientWorkloadExternalDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID:  "cw-12345",
				Name:        "Test CW UserGroup",
				Description: "Test description",
				IsActive:    true,
			},
		},
		Identities: []aembit.ClientWorkloadIdentityDTO{
			{
				Type:  "userGroup",
				Value: groupVal,
				Key:   idpID,
			},
		},
	}

	planModel := &models.ClientWorkloadResourceModel{}
	model := convertClientWorkloadDTOToModel(ctx, dto, planModel)

	var identities []models.IdentitiesModel
	diags := model.Identities.ElementsAs(ctx, &identities, false)
	if diags.HasError() {
		t.Fatalf("unexpected error extracting identities from model: %v", diags)
	}

	if len(identities) != 1 {
		t.Fatalf("expected 1 identity in model, got %d", len(identities))
	}

	if identities[0].Type.ValueString() != "userGroup" {
		t.Errorf("expected Type 'userGroup', got '%s'", identities[0].Type.ValueString())
	}
	if identities[0].Value.ValueString() != groupVal {
		t.Errorf("expected Value '%s', got '%s'", groupVal, identities[0].Value.ValueString())
	}
	if identities[0].IdentityProviderId.ValueString() != idpID {
		t.Errorf("expected IdentityProviderId '%s', got '%s'", idpID, identities[0].IdentityProviderId.ValueString())
	}
	if !identities[0].ClaimName.IsNull() {
		t.Errorf("expected ClaimName to be null, got '%s'", identities[0].ClaimName.ValueString())
	}
}

// TestValidateIdentities_UserGroupMissingIdpId verifies diagnostic error when identity_provider_id is missing.
func TestValidateIdentities_UserGroupMissingIdpId(t *testing.T) {
	t.Parallel()

	identities := []models.IdentitiesModel{
		{
			Type:               types.StringValue("userGroup"),
			Value:              types.StringValue("engineering"),
			IdentityProviderId: types.StringNull(),
			ClaimName:          types.StringNull(),
		},
	}

	diags := validateUserGroupIdentityTypeForConfig(identities)
	if !diags.HasError() {
		t.Fatalf("expected error for missing identity_provider_id, got none")
	}

	foundExpectedError := false
	for _, diag := range diags.Errors() {
		if strings.Contains(diag.Detail(), "Identity Provider ID is required when client identity type is userGroup.") {
			foundExpectedError = true
			break
		}
	}
	if !foundExpectedError {
		t.Errorf("did not find expected error message in: %v", diags.Errors())
	}

	// Also test empty string
	identitiesEmpty := []models.IdentitiesModel{
		{
			Type:               types.StringValue("userGroup"),
			Value:              types.StringValue("engineering"),
			IdentityProviderId: types.StringValue(""),
			ClaimName:          types.StringNull(),
		},
	}

	diagsEmpty := validateUserGroupIdentityTypeForConfig(identitiesEmpty)
	if !diagsEmpty.HasError() {
		t.Fatalf("expected error for empty identity_provider_id, got none")
	}

	// Test unknown string does not trigger error during ValidateConfig
	identitiesUnknown := []models.IdentitiesModel{
		{
			Type:               types.StringValue("userGroup"),
			Value:              types.StringValue("engineering"),
			IdentityProviderId: types.StringUnknown(),
			ClaimName:          types.StringNull(),
		},
	}

	diagsUnknown := validateUserGroupIdentityTypeForConfig(identitiesUnknown)
	if diagsUnknown.HasError() {
		t.Fatalf("expected no error for unknown identity_provider_id, got: %v", diagsUnknown.Errors())
	}
}

// TestValidateIdentities_UserGroupWithClaimName verifies diagnostic error when claim_name is set on userGroup.
func TestValidateIdentities_UserGroupWithClaimName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		claimName types.String
	}{
		{name: "non_empty_claim_name", claimName: types.StringValue("groups")},
		{name: "empty_claim_name", claimName: types.StringValue("")},
		{name: "unknown_claim_name", claimName: types.StringUnknown()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			identities := []models.IdentitiesModel{
				{
					Type:               types.StringValue("userGroup"),
					Value:              types.StringValue("engineering"),
					IdentityProviderId: types.StringValue("11111111-2222-3333-4444-555555555555"),
					ClaimName:          tc.claimName,
				},
			}

			diags := validateUserGroupIdentityTypeForConfig(identities)
			if !diags.HasError() {
				t.Fatalf("expected error for claim_name %v set on userGroup, got none", tc.claimName)
			}

			foundExpectedError := false
			for _, diag := range diags.Errors() {
				if strings.Contains(diag.Detail(), "Claim name is only applicable for oidcIdToken client identity type and cannot be specified for userGroup.") {
					foundExpectedError = true
					break
				}
			}
			if !foundExpectedError {
				t.Errorf("did not find expected error message in: %v", diags.Errors())
			}
		})
	}
}

// TestValidateIdentities_NonEmptyIdpIdOnNonGroupType verifies diagnostic error when identity_provider_id is set on k8sNamespace or oidcIdToken.
func TestValidateIdentities_NonEmptyIdpIdOnNonGroupType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		idType string
	}{
		{name: "k8sNamespace", idType: "k8sNamespace"},
		{name: "oidcIdToken", idType: "oidcIdToken"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			identities := []models.IdentitiesModel{
				{
					Type:               types.StringValue(tc.idType),
					Value:              types.StringValue("test-val"),
					IdentityProviderId: types.StringValue("11111111-2222-3333-4444-555555555555"),
					ClaimName:          types.StringNull(),
				},
			}

			diags := validateUserGroupIdentityTypeForConfig(identities)
			if !diags.HasError() {
				t.Fatalf("expected error for identity_provider_id on %s, got none", tc.idType)
			}

			foundExpectedError := false
			for _, diag := range diags.Errors() {
				if strings.Contains(diag.Detail(), "Identity Provider ID is only applicable for userGroup client identity type.") {
					foundExpectedError = true
					break
				}
			}
			if !foundExpectedError {
				t.Errorf("did not find expected error message in: %v", diags.Errors())
			}
		})
	}

	// Verify unknown identity_provider_id is also rejected for non-userGroup identity
	t.Run("unknown_identity_provider_id_rejected", func(t *testing.T) {
		identities := []models.IdentitiesModel{
			{
				Type:               types.StringValue("k8sNamespace"),
				Value:              types.StringValue("test-val"),
				IdentityProviderId: types.StringUnknown(),
				ClaimName:          types.StringNull(),
			},
		}

		diags := validateUserGroupIdentityTypeForConfig(identities)
		if !diags.HasError() {
			t.Fatalf("expected error for unknown identity_provider_id on k8sNamespace, got none")
		}

		foundExpectedError := false
		for _, diag := range diags.Errors() {
			if strings.Contains(diag.Detail(), "Identity Provider ID is only applicable for userGroup client identity type.") {
				foundExpectedError = true
				break
			}
		}
		if !foundExpectedError {
			t.Errorf("did not find expected error message in: %v", diags.Errors())
		}
	})
}

// TestValidateIdentities_MixedUserGroupAndRedirectUri verifies error when userGroup is mixed with oauthRedirectUri.
func TestValidateIdentities_MixedUserGroupAndRedirectUri(t *testing.T) {
	t.Parallel()

	identities := []models.IdentitiesModel{
		{
			Type:               types.StringValue("userGroup"),
			Value:              types.StringValue("engineering"),
			IdentityProviderId: types.StringValue("11111111-2222-3333-4444-555555555555"),
			ClaimName:          types.StringNull(),
		},
		{
			Type:               types.StringValue("oauthRedirectUri"),
			Value:              types.StringValue("https://example.com/oauth"),
			IdentityProviderId: types.StringNull(),
			ClaimName:          types.StringNull(),
		},
	}

	diags, hasMcp := validateMcpIdentityTypeForConfig(identities)
	if !hasMcp {
		t.Fatalf("expected hasMcp to be true")
	}
	if !diags.HasError() {
		t.Fatalf("expected error when mixing userGroup and oauthRedirectUri, got none")
	}

	foundExpectedError := false
	for _, diag := range diags.Errors() {
		if strings.Contains(diag.Detail(), "MCP client workload identities (`oauthRedirectUri` or `cimdClientId`) cannot be combined with other identity types.") {
			foundExpectedError = true
			break
		}
	}
	if !foundExpectedError {
		t.Errorf("did not find expected error message in: %v", diags.Errors())
	}
}
