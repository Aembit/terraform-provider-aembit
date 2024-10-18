package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const trustProviderPathRole string = "aembit_trust_provider.aws_role"
const trustProviderPathAzure string = "aembit_trust_provider.azure"
const trustProviderGitLab1 string = "aembit_trust_provider.gitlab1"
const trustProviderGitLab2 string = "aembit_trust_provider.gitlab2"

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
					resource.TestCheckResourceAttr(trustProviderPathAzure, "name", "TF Acceptance Azure"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderPathAzure, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderPathAzure, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: string(createFile), Check: testDeleteTrustProvider(trustProviderPathAzure), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: trustProviderPathAzure, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(trustProviderPathAzure, "name", "TF Acceptance Azure - Modified"),
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
					resource.TestCheckResourceAttr(trustProviderPathRole, "name", "TF Acceptance AWS Role"),
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
					resource.TestCheckResourceAttr(trustProviderPathRole, "name", "TF Acceptance AWS Role - Modified"),
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
					resource.TestCheckResourceAttr("aembit_trust_provider.aws", "name", "TF Acceptance AWS"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.aws", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.aws", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.aws", "name", "TF Acceptance AWS - Modified"),
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
					resource.TestCheckResourceAttr("aembit_trust_provider.gcp", "name", "TF Acceptance GCP Identity"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.gcp", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.gcp", "name", "TF Acceptance GCP Identity - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.gcp", "id"),
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
					resource.TestCheckResourceAttr("aembit_trust_provider.github", "name", "TF Acceptance GitHub Action"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.github", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.github", "name", "TF Acceptance GitHub Action - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.github", "id"),
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
					resource.TestCheckResourceAttr(trustProviderGitLab1, "name", "TF Acceptance GitLab Job1"),
					resource.TestCheckResourceAttr(trustProviderGitLab2, "name", "TF Acceptance GitLab Job2"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitLab1, "id"),
					resource.TestCheckResourceAttrSet(trustProviderGitLab2, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttr(trustProviderGitLab1, "gitlab_job.namespace_path", "namespace_path"),
					resource.TestCheckResourceAttr(trustProviderGitLab2, "gitlab_job.namespace_paths.0", "namespace_path1"),
					resource.TestCheckResourceAttr(trustProviderGitLab2, "gitlab_job.namespace_paths.1", "namespace_path2"),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderGitLab1, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(trustProviderGitLab1, "name", "TF Acceptance GitLab Job - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderGitLab1, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
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
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "name", "TF Acceptance Kerberos"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kerberos", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kerberos", "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsCount, "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsColor, "blue"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsDay, "Sunday"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.kerberos", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", "name", "TF Acceptance Kerberos - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsCount, "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsColor, "orange"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kerberos", tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTrustProviderResource_KubernetesServiceAccount(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/trust/kubernetes/TestAccTrustProviderResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "name", "TF Acceptance Kubernetes"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes", "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsCount, "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsColor, "blue"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsDay, "Sunday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes_key", "name", "TF Acceptance Kubernetes Key"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.kubernetes", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", "name", "TF Acceptance Kubernetes - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsCount, "2"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsColor, "orange"),
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes", tagsDay, "Tuesday"),
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr("aembit_trust_provider.kubernetes_key", "name", "TF Acceptance Kubernetes Key - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.kubernetes_key", "id"),
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
					resource.TestCheckResourceAttr("aembit_trust_provider.terraform", "name", "TF Acceptance Terraform Workspace"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
				),
			},
			// ImportState testing
			{ResourceName: "aembit_trust_provider.terraform", ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr("aembit_trust_provider.terraform", "name", "TF Acceptance Terraform Workspace - Modified"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_trust_provider.terraform", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
