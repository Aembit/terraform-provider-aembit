package provider

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	res "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestUnitAccessPolicyConfigure(t *testing.T) {
	var role accessPolicyResource = accessPolicyResource{}
	var configRequest res.ConfigureRequest = res.ConfigureRequest{ProviderData: "invalidData"}
	var configResponse res.ConfigureResponse = res.ConfigureResponse{}

	role.Configure(context.Background(), configRequest, &configResponse)

	assert.NotEmpty(t, configResponse.Diagnostics)
}

func TestAccAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_policy.first_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")

	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "clientworkloadNamespace", fmt.Sprintf("clientworkloadNamespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "aembit_access_policy.first_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttrSet("aembit_access_policy.first_policy", "id"),
					resource.TestCheckResourceAttrSet("aembit_access_policy.second_policy", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
