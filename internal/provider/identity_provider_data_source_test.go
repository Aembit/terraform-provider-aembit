package provider

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testIdpDataSource = "data.aembit_identity_providers.test_idps"

func TestAccIdentityProviderDataSource(t *testing.T) {
	createFile, err := os.ReadFile(
		"../../tests/identity_provider/data/TestAccIdentityProviderDataSource.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}

	startValue := "Identity Provider for TF Acceptance Test"
	createFileConfig, _, endValue := randomizeFileConfigs(string(createFile), "", startValue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testIdpDataSource, "identity_providers.#"),
					testAccCheckIdentityProviderExistsInDataSource(testIdpDataSource, endValue),
				),
			},
		},
	})
}

func testAccCheckIdentityProviderExistsInDataSource(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if ms.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		countStr := ms.Primary.Attributes["identity_providers.#"]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return err
		}

		for i := 0; i < count; i++ {
			if ms.Primary.Attributes[fmt.Sprintf("identity_providers.%d.name", i)] == name {
				return nil
			}
		}

		return fmt.Errorf("Identity Provider with name %s not found in %s", name, n)
	}
}
