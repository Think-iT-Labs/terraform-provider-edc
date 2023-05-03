package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyResource(t *testing.T) {
	resourceName := "edc_policy.pol"
	policyId := acctest.RandomWithPrefix("tf-acc-test")
	policyUID := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: testAccPolicyResourceConfig(policyId, policyUID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", policyId),
					resource.TestCheckResourceAttr(resourceName, "policy.uid", policyUID),
				),
			},
		},
	})
}

func testAccPolicyResourceConfig(policyId, policyUID string) string {
	return providerConfig + fmt.Sprintf(`
resource "edc_policy" "pol" {
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
`, policyId, policyUID)
}
