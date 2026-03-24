package provider

import (
	"context"
	"fmt"
	"os"
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
)

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

// TestAccClientWorkloadResource_k8sPodName_CustomResourceSet tests resource creation within a custom resource set.
func TestAccClientWorkloadResource_k8sPodName_CustomResourceSetAuth(t *testing.T) {
	skipNotCI(t)

	createFile, _ := os.ReadFile(
		"../../tests/client/resourceSet/TestAccClientWorkloadCustomResourceSet.tf",
	)
	createFileConfig, _, newName := randomizeFileConfigs(
		string(createFile),
		"",
		"custom-resource-set",
	)

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
				),
			},
			// ImportState testing
			{ResourceName: testCWResource, ImportState: true, ImportStateVerify: true},
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
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
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
						"https://test.aembit.local:12345",
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
				Config: string(modifyFile),
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

func TestConvertClientWorkloadDTOToModel_NormalizesSSOFields(t *testing.T) {
	t.Parallel()

	redirectModel := convertClientWorkloadDTOToModel(context.Background(), aembit.ClientWorkloadExternalDTO{
		EntityDTO: aembit.EntityDTO{
			ExternalID: "4fce7a89-9ab2-4fbf-a757-acbc8a1338d4",
			Name:       "redirect workload",
			IsActive:   true,
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
		EntityDTO: aembit.EntityDTO{
			ExternalID: "5fce7a89-9ab2-4fbf-a757-acbc8a1338d4",
			Name:       "non-redirect workload",
			IsActive:   true,
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
		values[i] = types.ObjectValueMust(models.TfIdentityObjectType.AttrTypes, map[string]attr.Value{
			"type":       identity.Type,
			"value":      identity.Value,
			"claim_name": identity.ClaimName,
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
