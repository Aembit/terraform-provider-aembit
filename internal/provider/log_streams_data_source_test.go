package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testLogStreamsDataSource string = "data.aembit_log_streams.test"
	testLogStreamResource    string = "aembit_log_stream.aws_s3_bucket"
)

func testFindLogStream(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetLogStream(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccLogStreamsDataSource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/log_stream/data/TestAccLogStreamsDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance AWSS3Bucket LogStream",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Integrations returned
					resource.TestCheckResourceAttrSet(testLogStreamsDataSource, "log_streams.#"),
					// Find newly created entry
					testFindLogStream(testLogStreamResource),
				),
			},
		},
	})
}
