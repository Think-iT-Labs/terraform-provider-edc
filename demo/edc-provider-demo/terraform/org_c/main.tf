resource "edc_asset" "asset_1" {
  provider = edc.org_c
  asset = {
    "asset:prop:id" : "asset_1_org_c",
    "asset:prop:name" : "FirstAssetOrgC",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_c1.txt"
    }
  }
}

resource "edc_asset" "asset_2" {
  provider = edc.org_c
  asset = {
    "asset:prop:id" : "asset_2_org_c",
    "asset:prop:name" : "SecondAssetOrgC",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_c2.txt"
    }
  }
}
