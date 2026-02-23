package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testTagsResource string = "aembit_client_workload.test_tags"
)

func TestAccTagsCreateBothTagsUpdateNoTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/tags/create_both_tags_update_no_tags/TestAccTags.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_both_tags_update_no_tags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"4",
					),
					resource.TestCheckResourceAttr(testTagsResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testTagsResource, tagsDay, "Sunday"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllName,
						"Terraform",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllOwner,
						"Aembit",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagsCreateDefaultTagsUpdateNoTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/tags/create_defaulttags_update_no_tags/TestAccTags.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_defaulttags_update_no_tags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllName,
						"Terraform",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllOwner,
						"Aembit",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagsCreateNoTagsUpdateBothTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_both_tags/TestAccTags.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_both_tags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"4",
					),
					resource.TestCheckResourceAttr(testTagsResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testTagsResource, tagsDay, "Sunday"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllName,
						"Terraform",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllOwner,
						"Aembit",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagsCreateNoTagsUpdateDefaultTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_defaulttags/TestAccTags.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_defaulttags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"2",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllName,
						"Terraform",
					),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllOwner,
						"Aembit",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagsCreateNoTagsUpdateResourceTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_resourcetags/TestAccTags.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_no_tags_update_resourcetags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"2",
					),
					resource.TestCheckResourceAttr(testTagsResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testTagsResource, tagsDay, "Sunday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccTagsCreateResourceTagsUpdateNoTags(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/tags/create_resourcetags_update_no_tags/TestAccTags.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/tags/create_resourcetags_update_no_tags/TestAccTags.tfmod",
	)
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
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
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(
						testTagsResource,
						tagsAllCount,
						"2",
					),
					resource.TestCheckResourceAttr(testTagsResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testTagsResource, tagsDay, "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testTagsResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testTagsResource, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Tags.
					resource.TestCheckResourceAttr(testTagsResource, tagsCount, "0"),
					resource.TestCheckResourceAttr(testTagsResource, tagsAllCount, "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
