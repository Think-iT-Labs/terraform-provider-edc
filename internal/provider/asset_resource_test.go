package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: testAccAssetResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edc_asset.s3", "id", "test_id"),
					resource.TestCheckResourceAttr("edc_asset.s3.asset", "asset:prop:name", "test asset"),
					resource.TestCheckResourceAttr("edc_asset.s3.asset", "asset:prop:contenttype", "application/json"),
					resource.TestCheckResourceAttr("edc_asset.s3.asset", "asset:prop:id", "test-asset-id"),
					resource.TestCheckResourceAttr("edc_asset.s3.data.s3", "type", "AmazonS3"),
					resource.TestCheckResourceAttr("edc_asset.s3.data.s3", "name", "test file"),
					resource.TestCheckResourceAttr("edc_asset.s3.data.s3", "bucket_name", "testBucket"),
					resource.TestCheckResourceAttr("edc_asset.s3.data.s3", "access_key_id", "dummy_key"),
					resource.TestCheckResourceAttr("edc_asset.s3.data.s3", "secret_access_key", "dummy_key"),
				),
			},
		},
	})
}

func testAccAssetResourceConfig() string {
	return `
resource "edc_asset" "s3" {
	asset = {
		"asset:prop:name" : "test asset",
		"asset:prop:contenttype" : "application/json",
		"asset:prop:id": "test-asset-id",
	}

	data = {
		s3 = {
		type              = "AmazonS3"
		name              = "test file"
		bucket_name       = "testBucket"
		access_key_id     = "dummy_key"
		secret_access_key = "dummy_key"
		}
	}
	}
`
}
