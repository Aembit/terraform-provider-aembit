package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	trustProviderPathRole          = "aembit_trust_provider.aws_role"
	trustProviderPathAzure         = "aembit_trust_provider.azure"
	trustProviderGitLab1           = "aembit_trust_provider.gitlab1"
	trustProviderGitLab2           = "aembit_trust_provider.gitlab2"
	trustProviderGitLabMixed       = "aembit_trust_provider.gitlab_mixed"
	gitLabOidcClientID             = "gitlab_job.oidc_client_id"
	gitLabIdentityArnMatch         = ":identity:gitlab_idtoken:"
	trustProviderAwsPath           = "aembit_trust_provider.aws"
	trustProviderGcp               = "aembit_trust_provider.gcp"
	trustProviderGitHub            = "aembit_trust_provider.github"
	trustProviderKerberos          = "aembit_trust_provider.kerberos"
	trustProviderTerraformResource = "aembit_trust_provider.terraform"
)

func testDeleteTrustProvider(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteTrustProvider(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccTrustProviderResource_AzureMetadata(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/azure/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/azure/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderPathAzure,
						"name",
						"TF Acceptance Azure",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderPathAzure, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderPathAzure, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteTrustProvider(trustProviderPathAzure),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: trustProviderPathAzure, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderPathAzure,
						"name",
						"TF Acceptance Azure - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_AwsRole(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/aws_role/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/aws_role/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderPathRole,
						"name",
						"TF Acceptance AWS Role",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderPathRole, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderPathRole, "id"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderPathRole, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderPathRole,
						"name",
						"TF Acceptance AWS Role - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderPathRole, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_AwsMetadata(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/aws/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/aws/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderAwsPath,
						"name",
						"TF Acceptance AWS",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderAwsPath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderAwsPath, "id"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderAwsPath, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderAwsPath,
						"name",
						"TF Acceptance AWS - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GcpIdentity(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/gcp/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/gcp/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderGcp,
						"name",
						"TF Acceptance GCP Identity",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGcp, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderGcp, "id"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderGcp, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderGcp,
						"name",
						"TF Acceptance GCP Identity - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGcp, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GitHubAction(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/github/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/github/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"name",
						"TF Acceptance GitHub Action",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderGitHub,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"name",
						"TF Acceptance GitHub Action - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GitHubAction_OidcEndpoint(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/trust/github/oidc_endpoint/TestAccTrustProviderResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/trust/github/oidc_endpoint/TestAccTrustProviderResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"name",
						"TF Acceptance GitHub Action",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"github_action.oidc_endpoint",
						"https://gitlab.com",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderGitHub,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"name",
						"TF Acceptance GitHub Action - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitHub, "id"),
					resource.TestCheckResourceAttr(
						trustProviderGitHub,
						"github_action.oidc_endpoint",
						"https://token.actions.githubusercontent.com",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GitLabJob(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/gitlab/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/gitlab/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderGitLab1,
						"name",
						"TF Acceptance GitLab Job1",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLab2,
						"name",
						"TF Acceptance GitLab Job2",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLabMixed,
						"name",
						"TF Acceptance GitLab Mixed",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitLab1, "id"),
					resource.TestCheckResourceAttrSet(trustProviderGitLab2, "id"),
					resource.TestCheckResourceAttrSet(trustProviderGitLabMixed, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttr(
						trustProviderGitLab1,
						"gitlab_job.oidc_endpoint",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLab2,
						"gitlab_job.oidc_endpoint",
						"https://gitlab.com",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLab1,
						"gitlab_job.namespace_path",
						"namespace_path",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLab2,
						"gitlab_job.namespace_paths.0",
						"namespace_path1",
					),
					resource.TestCheckResourceAttr(
						trustProviderGitLab2,
						"gitlab_job.namespace_paths.1",
						"namespace_path2",
					),
					// Check read-only values
					checkValidClientID(
						trustProviderGitLab1,
						gitLabOidcClientID,
						gitLabIdentityArnMatch,
					),
					checkValidClientID(
						trustProviderGitLab2,
						gitLabOidcClientID,
						gitLabIdentityArnMatch,
					),
					checkValidClientID(
						trustProviderGitLabMixed,
						gitLabOidcClientID,
						gitLabIdentityArnMatch,
					),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderGitLab1, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderGitLab1,
						"name",
						"TF Acceptance GitLab Job - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitLab1, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_GitLabJob_Validation(t *testing.T) {
	invalidNameFile, _ := os.ReadFile(
		"../../tests/trust/gitlab/TestAccTrustProviderResource.tfinvalid",
	)

	regexChecks := []string{
		// Protect against empty strings
		`Attribute gitlab_job.namespace_path string length must be at least 1, got: 0`,
		`Attribute gitlab_job.ref_path string length must be at least 1, got: 0`,
		`Attribute gitlab_job.project_path string length must be at least 1, got: 0`,
		`Attribute gitlab_job.subject string length must be at least 1, got: 0`,
		// Protect against sets with fewer than 2 items
		`Attribute gitlab_job.namespace_paths set must contain at least 2 elements`,
		`Attribute gitlab_job.ref_paths set must contain at least 2 elements`,
		`Attribute gitlab_job.project_paths set must contain at least 2 elements`,
		`Attribute gitlab_job.subjects set must contain at least 2 elements`,
		// Protect against sets with empty strings
		`(?s)Error: Invalid Attribute Value Length(.*)with aembit_trust_provider.gitlab_set_strings,`,
		// Protect against conflicts
		`(?s)These attributes cannot be configured together:(.{2})[gitlab_job.namespace_path,gitlab_job.namespace_paths]`,
		`(?s)These attributes cannot be configured together:(.{2})[gitlab_job.ref_path,gitlab_job.ref_paths]`,
		`(?s)These attributes cannot be configured together:(.{2})[gitlab_job.project_path,gitlab_job.project_paths]`,
		`(?s)These attributes cannot be configured together:(.{2})[gitlab_job.subject,gitlab_job.subjects]`,
	}
	validationChecks := []resource.TestStep{}
	for _, check := range regexChecks {
		validationChecks = append(
			validationChecks,
			resource.TestStep{
				Config:      string(invalidNameFile),
				ExpectError: regexp.MustCompile(check),
			},
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    validationChecks,
	})
}

func TestAccTrustProviderResource_Kerberos(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/kerberos/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/kerberos/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						"name",
						"TF Acceptance Kerberos",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKerberos, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKerberos, "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsColor,
						"blue",
					),
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsDay,
						"Sunday",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderKerberos,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						"name",
						"TF Acceptance Kerberos - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsColor,
						"orange",
					),
					resource.TestCheckResourceAttr(
						trustProviderKerberos,
						tagsDay,
						"Tuesday",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_KubernetesServiceAccount(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tfmod")

	const trustProviderKubernetes string = "aembit_trust_provider.kubernetes"
	const trustProviderKubernetesKey string = "aembit_trust_provider.kubernetes_key"
	const trustProviderKubernetesJWKS string = "aembit_trust_provider.kubernetes_jwks"
	const trustProviderKubernetesSymmetricKey string = "aembit_trust_provider.kubernetes_symmetric_key"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetes,
						"name",
						"TF Acceptance Kubernetes",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetes, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetes, "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsCount, "2"),
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsColor, "blue"),
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsDay, "Sunday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesKey,
						"name",
						"TF Acceptance Kubernetes Key",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesKey, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesJWKS,
						"name",
						"TF Acceptance Kubernetes JWKS",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesJWKS, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesJWKS, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesSymmetricKey,
						"name",
						"TF Acceptance Kubernetes Symmetric Key",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesSymmetricKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesSymmetricKey, "id"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderKubernetes, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderKubernetes,
						"name",
						"TF Acceptance Kubernetes - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsCount, "2"),
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsColor, "orange"),
					resource.TestCheckResourceAttr(trustProviderKubernetes, tagsDay, "Tuesday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesKey,
						"name",
						"TF Acceptance Kubernetes Key - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesKey, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesJWKS,
						"name",
						"TF Acceptance Kubernetes JWKS - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesJWKS, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesJWKS, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderKubernetesSymmetricKey,
						"name",
						"TF Acceptance Kubernetes Symmetric Key - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderKubernetesSymmetricKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderKubernetesSymmetricKey, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_TerraformWorkspace(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/terraform/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/terraform/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderTerraformResource,
						"name",
						"TF Acceptance Terraform Workspace",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderTerraformResource, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderTerraformResource, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      trustProviderTerraformResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderTerraformResource,
						"name",
						"TF Acceptance Terraform Workspace - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderTerraformResource, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_OidcIdToken(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/oidc-id-token/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/trust/oidc-id-token/TestAccTrustProviderResource.tfmod",
	)

	const trustProviderOidcidToken = "aembit_trust_provider.oidcidtoken"
	const trustProviderOidcidTokenKey = "aembit_trust_provider.oidcidtoken_key"
	const trustProviderOidcidTokenJWKS = "aembit_trust_provider.oidcidtoken_jwks"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderOidcidToken,
						"name",
						"TF Acceptance OIDC ID Token",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderOidcidToken, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderOidcidToken, "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsCount, "2"),
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsColor, "blue"),
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsDay, "Sunday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderOidcidTokenKey,
						"name",
						"TF Acceptance OIDC ID Token Key",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenKey, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderOidcidTokenJWKS,
						"name",
						"TF Acceptance OIDC ID Token JWKS",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenJWKS, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenJWKS, "id"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderOidcidToken, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						trustProviderOidcidToken,
						"name",
						"TF Acceptance OIDC ID Token - Modified",
					),
					// Verify Tags.
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsCount, "2"),
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsColor, "blue"),
					resource.TestCheckResourceAttr(trustProviderOidcidToken, tagsDay, "Sunday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderOidcidTokenKey,
						"name",
						"TF Acceptance OIDC ID Token Key - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenKey, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenKey, "id"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(
						trustProviderOidcidTokenJWKS,
						"name",
						"TF Acceptance OIDC ID Token JWKS - Modified",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenJWKS, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderOidcidTokenJWKS, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_OidcIdToken_MissingRequiredRSAField(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/trust/oidc-id-token/TestAccTrustProviderResource_MissingRequiredRSAField.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      string(createFile),
				ExpectError: regexp.MustCompile(`does not have RSA required fields: e, n`),
			},
		},
	})
}

func TestAccTrustProviderResource_OidcIdToken_MissingRequiredECDSAField(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/trust/oidc-id-token/TestAccTrustProviderResource_MissingRequiredEDSAField.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      string(createFile),
				ExpectError: regexp.MustCompile(`does not have ECDSA required fields: x, y, crv`),
			},
		},
	})
}
