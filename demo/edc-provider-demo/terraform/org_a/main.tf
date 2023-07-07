resource "edc_asset" "asset_1" {
  asset = {
    "asset:prop:id" : "asset_1_org_a",
    "asset:prop:name" : "FirstAssetOrgA",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    http = {
      base_url = "http://nginx:80/file_a1.txt"
      method   = "GET"
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
      base_url = "http://nginx:80/file_a2.txt"
      method   = "GET"
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
      base_url = "http://nginx:80/file_a3.txt"
      method   = "GET"
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
      base_url = "http://nginx:80/file_a4.txt"
      method   = "GET"
    }
  }
}

resource "edc_policy" "policy" {
  id = "abcdPolicy"
  policy = {
    permissions = [
      {
        edctype = "dataspaceconnector:permission",
        target  = edc_asset.asset_1.id,
        action = {
          type = "USE"
        },
        type = "set"
      }
    ]
  }
}

resource "edc_contract_definition" "name" {
  access_policy_id   = edc_policy.policy.id
  contract_policy_id = edc_policy.policy.id
  validity           = 31536000
  criteria = [
    {
      operand_left  = "asset:prop:id"
      operator      = "="
      operand_right = edc_asset.asset_1.id
    }
  ]
}
