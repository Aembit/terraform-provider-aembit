package provider

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testCallerIdentityDataSource = "data.aembit_caller_identity.test"
)

// Basic check that tenant_id is present in state
func testCallerIdentityTenantID(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if tenantID, ok := rs.Primary.Attributes["tenant_id"]; !ok || tenantID == "" {
			return fmt.Errorf("expected tenant_id to be set, got: %v", tenantID)
		}
		return nil
	}
}

func TestAccCallerIdentityDataSource(t *testing.T) {
	createFile, err := os.ReadFile("../../tests/caller_identity/data/TestAccCallerIdentityDataSource.tf")
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}

	config := string(createFile)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCallerIdentityDataSource, "tenant_id"),
					testCallerIdentityTenantID(testCallerIdentityDataSource),
				),
			},
		},
	})
}
