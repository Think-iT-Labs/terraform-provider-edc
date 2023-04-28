package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetResource(t *testing.T) {
	resourceName := "edc_asset.s3"
	assetId := "test-asset-id"
	assetName := "test asset"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: testAccAssetResourceConfig(assetId, assetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", assetId),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:name", assetName),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:contenttype", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:id", assetId),
					resource.TestCheckResourceAttr(resourceName, "data.s3.name", "test file"),
					resource.TestCheckResourceAttr(resourceName, "data.s3.bucket_name", "testBucket"),
					resource.TestCheckResourceAttr(resourceName, "data.s3.access_key_id", "dummy_key"),
					resource.TestCheckResourceAttr(resourceName, "data.s3.secret_access_key", "dummy_key"),
				),
			},
		},
	})
}

func testAccAssetResourceConfig(assetId, assetName string) string {
	return providerConfig + fmt.Sprintf(`
resource "edc_asset" "s3" {
	asset = {
		"asset:prop:name" : %[2]q,
		"asset:prop:contenttype" : "application/json",
		"asset:prop:id": %[1]q,
	}

	data = {
		s3 = {
			name              = "test file"
			bucket_name       = "testBucket"
			access_key_id     = "dummy_key"
			secret_access_key = "dummy_key"
		}
	}
}
`, assetId, assetName)
}
