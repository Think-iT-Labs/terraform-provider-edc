resource "edc_asset" "asset_1" {
  asset = {
    "asset:prop:id" : "asset_1_org_a",
    "asset:prop:name" : "FirstAssetOrgA",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_a1.txt"
    }
  }
}

resource "edc_asset" "asset_2" {
  asset = {
    "asset:prop:id" : "asset_2_org_a",
    "asset:prop:name" : "SecondAssetOrgA",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_a2.txt"
    }
  }
}

resource "edc_asset" "asset_3" {
  asset = {
    "asset:prop:id" : "asset_3_org_a",
    "asset:prop:name" : "ThirdAssetOrgA",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_a3.txt"
    }
  }
}

resource "edc_asset" "asset_4" {
  asset = {
    "asset:prop:id" : "asset_4_org_a",
    "asset:prop:name" : "FourthAssetOrgA",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://localhost:8080"
      method   = "GET"
      path     = "/file_a4.txt"
    }
  }
}
