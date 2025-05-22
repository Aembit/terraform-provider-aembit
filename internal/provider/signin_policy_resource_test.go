package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testSignInPolicy string = "aembit_signin_policy.test"

func TestAccSigninPolicy(t *testing.T) {
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

func TestAccSigninPolicy_MissingRequiredAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aembit_signin_policy" "test" {
  // Missing required attributes: sso_required, mfa_required
}
`,
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccSigninPolicy_InvalidAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `
resource "aembit_signin_policy" "test" {
    sso_required = "invalid_bool" // Should be a boolean
  	mfa_required  = true
}
`,
				ExpectError: regexp.MustCompile(`Inappropriate value for attribute`),
			},
		},
	})
}
