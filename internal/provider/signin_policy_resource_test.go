package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testSignInPolicy string = "aembit_signin_policy.test"

func TestAccSigninPolicy(t *testing.T) {
	tfVersion := getTerraformVersion()
	if tfVersion != "v1.9" {
		t.Skip("Skipping testing in Terraform version " + tfVersion)
	}

	createFile, _ := os.ReadFile("../../tests/signin_policy/TestAccSignInPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/signin_policy/TestAccSignInPolicyResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testSignInPolicy, "mfa_required", "true"),
					resource.TestCheckResourceAttr(testSignInPolicy, "sso_required", "true"),
				),
				ResourceName: testSignInPolicy,
			},
			// ImportState testing
			{ResourceName: testSignInPolicy, ImportState: true, ImportStateVerify: true},
			// Update
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(testSignInPolicy, "mfa_required", "true"),
					resource.TestCheckResourceAttr(testSignInPolicy, "sso_required", "false"),
				),
				ResourceName: testSignInPolicy,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
