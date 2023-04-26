package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAssetDataSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.edc_asset.my_asset", "id", "assetId"),
					resource.TestCheckResourceAttr("data.edc_asset.my_asset", "createdAt", ""),
					resource.TestCheckResourceAttr("data.edc_asset.my_asset", "asset_properties.asset:prop:name", "1"),
					resource.TestCheckResourceAttr("data.edc_asset.my_asset", "asset_properties.asset:prop:contenttype", "application/json"),
					resource.TestCheckResourceAttr("data.edc_asset.my_asset", "asset_properties.asset:prop:id", "assetId"),
				),
			},
		},
	})
}

func testAccAssetDataSource() string {
	return providerConfig + `data "edc_asset" "my_asset" {
		id: "assetId"
	}`
}
