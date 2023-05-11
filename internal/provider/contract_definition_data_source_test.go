package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestContractDefinitionDataSource(t *testing.T) {
	resourceName := "edc_contract_definition.test"
	dataSourceName := "data.edc_contract_definition.my_contract"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccContractDefinitionDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccContractDefinitionDataSourceConfig() string {
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

data "edc_contract_definition" "my_contract" {
	id  = edc_contract_definition.test.id
}
`
}
