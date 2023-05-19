package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccS3AssetResource(t *testing.T) {
	resourceName := "edc_asset.s3"
	assetId := acctest.RandomWithPrefix("tf-acc-test")
	assetName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: testAccS3AssetResourceConfig(assetId, assetName),
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

func testAccS3AssetResourceConfig(assetId, assetName string) string {
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

func TestAccHttpAssetResource(t *testing.T) {
	resourceName := "edc_asset.http"
	assetId := acctest.RandomWithPrefix("tf-acc-test")
	assetName := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: testAccHttpAssetResourceConfig(assetId, assetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", assetId),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:name", assetName),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:contenttype", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "asset.asset:prop:id", assetId),
					resource.TestCheckResourceAttr(resourceName, "data.http.name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "data.http.base_url", "https://connecor/invalid-data.json"),
				),
			},
		},
	})
}

func testAccHttpAssetResourceConfig(assetId, assetName string) string {
	return providerConfig + fmt.Sprintf(`
resource "edc_asset" "http" {
	asset = {
		"asset:prop:name" : %[2]q,
		"asset:prop:contenttype" : "application/json",
		"asset:prop:id": %[1]q,
	}

	data = {
		http = {
		  name  = "terraform"
		  base_url = "https://connecor/invalid-data.json"
		}
	}
}
`, assetId, assetName)
}
