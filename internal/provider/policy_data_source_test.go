package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyDataSource(t *testing.T) {
	dataSourceName := "data.edc_policy.policy"
	resourceName := "edc_policy.policy"

	policyId := acctest.RandomWithPrefix("tf-acc-test")
	policyUID := acctest.RandomWithPrefix("tf-acc-test")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccPolicyDataSource(policyId, policyUID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "createdAt", dataSourceName, "createdAt"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccPolicyDataSource(policyId, policyUid string) string {
	return providerConfig + fmt.Sprintf(`
resource "edc_policy" "policy" {
	id = %[1]q
	policy = {
		uid = %[2]q,
		permissions = [
			{
				target = "assetId",
				action = {
					type = "USE"
				},
				edctype = "dataspaceconnector:permission"
			}
		]
	}
}

data "edc_policy" "policy" {
	id = edc_policy.policy.id
}
`, policyId, policyUid)
}
