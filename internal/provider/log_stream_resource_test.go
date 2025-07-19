package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testLogStreamAWSS3Bucket                   string = "aembit_log_stream.aws_s3_bucket"
	testLogStreamGCSBucket                     string = "aembit_log_stream.gcs_bucket"
	testLogStreamSplunkHttpEventCollector      string = "aembit_log_stream.splunk_http_event_collector"
	testLogStreamCrowdstrikeHttpEventCollector string = "aembit_log_stream.crowdstrike_http_event_collector"
)

func testDeleteLogStream(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteLogStream(context.Background(), rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccLogStreamResource_AWSS3Bucket(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/log_stream/awsS3Bucket/TestAccLogStreamResource.tf")
	modifyFile, _ := os.ReadFile(
		"../../tests/log_stream/awsS3Bucket/TestAccLogStreamResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify LogStream Name
					resource.TestCheckResourceAttr(
						testLogStreamAWSS3Bucket,
						"name",
						"TF Acceptance AWSS3Bucket LogStream",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testLogStreamAWSS3Bucket, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testLogStreamAWSS3Bucket, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteLogStream(testLogStreamAWSS3Bucket),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testLogStreamAWSS3Bucket, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testLogStreamAWSS3Bucket,
						"name",
						"TF Acceptance AWSS3Bucket LogStream - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccLogStreamResource_GCSBucket(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/log_stream/gcsBucket/TestAccLogStreamResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/log_stream/gcsBucket/TestAccLogStreamResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify LogStream Name
					resource.TestCheckResourceAttr(
						testLogStreamGCSBucket,
						"name",
						"TF Acceptance GCSBucket LogStream",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testLogStreamGCSBucket, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testLogStreamGCSBucket, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteLogStream(testLogStreamGCSBucket),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testLogStreamGCSBucket, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testLogStreamGCSBucket,
						"name",
						"TF Acceptance GCSBucket LogStream - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccLogStreamResource_Splunk(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/log_stream/splunkHttpEventCollector/TestAccLogStreamResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/log_stream/splunkHttpEventCollector/TestAccLogStreamResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify LogStream Name
					resource.TestCheckResourceAttr(
						testLogStreamSplunkHttpEventCollector,
						"name",
						"TF Acceptance SplunkHttpEventCollector LogStream",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testLogStreamSplunkHttpEventCollector, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testLogStreamSplunkHttpEventCollector, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteLogStream(testLogStreamSplunkHttpEventCollector),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{
				ResourceName:      testLogStreamSplunkHttpEventCollector,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testLogStreamSplunkHttpEventCollector,
						"name",
						"TF Acceptance SplunkHttpEventCollector LogStream - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccLogStreamResource_Crowdstrike(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/log_stream/crowdstrikeHttpEventCollector/TestAccLogStreamResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/log_stream/crowdstrikeHttpEventCollector/TestAccLogStreamResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify LogStream Name
					resource.TestCheckResourceAttr(
						testLogStreamCrowdstrikeHttpEventCollector,
						"name",
						"TF Acceptance CrowdstrikeHttpEventCollector LogStream",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(
						testLogStreamCrowdstrikeHttpEventCollector,
						"id",
					),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(
						testLogStreamCrowdstrikeHttpEventCollector,
						"id",
					),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{
				Config:             string(createFile),
				Check:              testDeleteLogStream(testLogStreamCrowdstrikeHttpEventCollector),
				ExpectNonEmptyPlan: true,
			},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{
				ResourceName:      testLogStreamCrowdstrikeHttpEventCollector,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(
						testLogStreamCrowdstrikeHttpEventCollector,
						"name",
						"TF Acceptance CrowdstrikeHttpEventCollector LogStream - Modified",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
