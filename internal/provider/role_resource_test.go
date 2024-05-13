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

func TestUnitRoleConfigure(t *testing.T) {
	var role roleResource = roleResource{}
	var configRequest res.ConfigureRequest = res.ConfigureRequest{ProviderData: "invalidData"}
	var configResponse res.ConfigureResponse = res.ConfigureResponse{}

	role.Configure(context.Background(), configRequest, &configResponse)

	assert.NotEmpty(t, configResponse.Diagnostics)
}

func TestAccRoleResource(t *testing.T) {
	const resourceID string = "aembit_role.role"
	createFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tfmod")

	randID := rand.Intn(10000000)
	createResourceName := fmt.Sprintf("TF Acceptance Role %d", randID)
	modifyResourceName := fmt.Sprintf("TF Acceptance Role %d - Modified", randID)
	createFileConfig := strings.ReplaceAll(string(createFile), "TF Acceptance Role", createResourceName)
	modifyFileConfig := strings.ReplaceAll(string(modifyFile), "TF Acceptance Role - Modified", modifyResourceName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(resourceID, "name", createResourceName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(resourceID, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(resourceID, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(resourceID, "name", modifyResourceName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
