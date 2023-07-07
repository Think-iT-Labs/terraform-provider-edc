resource "edc_asset" "asset_1" {
  asset = {
    "asset:prop:id" : "asset_1_org_b",
    "asset:prop:name" : "FirstAssetOrgB",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080/file_b1.txt"
      method   = "GET"
    }
  }
}
