package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetDataSource(t *testing.T) {
	dataSourceName := "data.edc_asset.my_asset"
	resourceName := "edc_asset.s3"

	assetName := acctest.RandomWithPrefix("tf-acc-test")
	assetId := acctest.RandomWithPrefix("tf-acc-test")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAssetDataSource(assetId, assetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "createdAt", dataSourceName, "createdAt"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "asset.asset:prop:name", dataSourceName, "asset_properties.asset:prop:name"),
					resource.TestCheckResourceAttrPair(resourceName, "asset.asset:prop:contenttype", dataSourceName, "asset_properties.asset:prop:contenttype"),
					resource.TestCheckResourceAttrPair(resourceName, "data.s3.name", dataSourceName, "data_address.name"),
					resource.TestCheckResourceAttrPair(resourceName, "data.s3.bucket_name", dataSourceName, "data_address.bucketName"),
					resource.TestCheckResourceAttrPair(resourceName, "data.s3.access_key_id", dataSourceName, "data_address.accessKeyId"),
					resource.TestCheckResourceAttrPair(resourceName, "data.s3.secret_access_key", dataSourceName, "data_address.secretAccessKey"),
				),
			},
		},
	})
}

func testAccAssetDataSource(assetId, assetName string) string {
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
data "edc_asset" "my_asset" {
	id = edc_asset.s3.id
}
`, assetId, assetName)
}
