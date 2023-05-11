package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContractDefinitionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// Steps: []resource.TestStep{
		// 		}
	})
}

func testAccContractDefinitionResourceConfig() string {
	return providerConfig + fmt.Sprint(`
resource "contract_definition" "test" {
	id = "test"
	access_policy_id = "test"
	contract_policy_id = "test"
	validity = 600
	criteria = [
		{
			operator = "eq"
			operand_left = "test"
			operand_right = "test"
		}
	]
}
`)
}
