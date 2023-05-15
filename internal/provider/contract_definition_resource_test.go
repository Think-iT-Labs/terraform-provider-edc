package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContractDefinitionResource(t *testing.T) {
	resourceName := "edc_contract_definition.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractDefinitionResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "access_policy_id", "test"),
					resource.TestCheckResourceAttr(resourceName, "contract_policy_id", "test"),
					resource.TestCheckResourceAttr(resourceName, "validity", "600"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operator", "eq"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operand_left", "test"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operand_right", "test"),
				),
			},
		},
	})
}

func testAccContractDefinitionResourceConfig() string {
	return providerConfig + `
resource "edc_contract_definition" "test" {
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
`
}
