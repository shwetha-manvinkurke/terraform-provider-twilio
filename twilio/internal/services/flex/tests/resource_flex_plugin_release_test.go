package tests

import (
	"fmt"
	"testing"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/common"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/internal/acceptance"
	"github.com/RJPearson94/terraform-provider-twilio/twilio/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var pluginReleaseResourceName = "twilio_flex_plugin_release"

func TestAccTwilioFlexPluginRelease_basic(t *testing.T) {
	stateResourceName := fmt.Sprintf("%s.plugin_release", pluginReleaseResourceName)

	name := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTwilioFlexPluginReleaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTwilioFlexPluginRelease_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTwilioFlexPluginReleaseExists(stateResourceName),
					resource.TestCheckResourceAttrSet(stateResourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "configuration_sid"),
					resource.TestCheckResourceAttrSet(stateResourceName, "date_created"),
					resource.TestCheckResourceAttrSet(stateResourceName, "url"),
				),
			},
			{
				ResourceName:      stateResourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTwilioFlexPluginReleaseImportStateIdFunc(stateResourceName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTwilioFlexPluginReleaseDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Flex

	for _, rs := range s.RootModule().Resources {
		if rs.Type != pluginReleaseResourceName {
			continue
		}

		if _, err := client.PluginRelease(rs.Primary.ID).Fetch(); err != nil {
			if utils.IsNotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error occurred when retrieving plugin release information %s", err.Error())
		}
	}

	return nil
}

func testAccCheckTwilioFlexPluginReleaseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*common.TwilioClient).Flex

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if _, err := client.PluginRelease(rs.Primary.ID).Fetch(); err != nil {
			return fmt.Errorf("Error occurred when retrieving plugin release information %s", err.Error())
		}

		return nil
	}
}

func testAccTwilioFlexPluginReleaseImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Not found: %s", name)
		}

		return fmt.Sprintf("/PluginService/Releases/%s", rs.Primary.Attributes["sid"]), nil
	}
}

func testAccTwilioFlexPluginRelease_basic(name string) string {
	return fmt.Sprintf(`
resource "twilio_flex_plugin_configuration" "plugin_configuration" {
  name = "%s"
}

resource "twilio_flex_plugin_release" "plugin_release" {
  configuration_sid = twilio_flex_plugin_configuration.plugin_configuration.sid
}
`, name)
}
