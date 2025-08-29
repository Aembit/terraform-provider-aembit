package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testCallerIdentityDataSource          = "data.aembit_caller_identity.test"
	testCallerIdentitySecondaryDataSource = "data.aembit_caller_identity.secondary"
)

func TestAccCallerIdentityDataSource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/caller_identity/data/TestAccCallerIdentityDataSource.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						testCallerIdentityDataSource,
						"tenant_id",
					),
				),
			},
		},
	})
}

func TestAccCallerIdentityDataSource_SecondaryProvider(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/caller_identity/data/TestAccCallerIdentityDataSourceSecondaryProvider.tf",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						testCallerIdentitySecondaryDataSource,
						"tenant_id",
					),
				),
			},
		},
	})
}
