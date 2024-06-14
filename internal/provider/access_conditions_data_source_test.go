package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccessConditionsDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/data/TestAccAccessConditionsDataSource.tf")

	//randID := rand.Intn(10000000)
	//createFileConfig := strings.ReplaceAll(string(createFile), "unittest1namespace", fmt.Sprintf("unittest1namespace%d", randID))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Access Conditions returned
					resource.TestCheckResourceAttrSet("data.aembit_access_conditions.test", "access_conditions.#"),
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.name", "TF Acceptance Crowdstrike"),
					// Verify Tags.
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.%", "2"),
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.color", "blue"),
					resource.TestCheckResourceAttr("data.aembit_access_conditions.test", "access_conditions.0.tags.day", "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_access_conditions.test", "access_conditions.0.id"),
				),
			},
		},
	})
}

/*
func testAccCheckExampleWidgetValues(widget *example.Widget, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *widget.Active != true {
			return fmt.Errorf("bad active state, expected \"true\", got: %#v", *widget.Active)
		}
		if *widget.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, *widget.Name)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckAccessConditionExists(accessConditionID string, accessCondition *aembit.AccessConditionDTO) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[accessConditionID]
		if !ok {
			return fmt.Errorf("Access Condition ID not found: %s", accessConditionID)
		}

		// retrieve the configured client from the test setup
		conn := testAccProvider.Meta().(*aembit.CloudClient)
		resp, err := conn.DescribeWidget(&example.DescribeWidgetsInput{
			WidgetIdentifier: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		if resp.Widget == nil {
			return fmt.Errorf("Widget (%s) not found", rs.Primary.ID)
		}

		// assign the response Widget attribute to the widget pointer
		*widget = *resp.Widget

		return nil
	}
}
*/
