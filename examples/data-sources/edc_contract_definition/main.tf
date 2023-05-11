terraform {
  required_providers {
    edc = {
      source = "Think-iT-Labs/edc"
    }
  }
}

provider "edc" {
  token = "1234"
  addresses = {
    default    = "http://localhost:29193/api"
    management = "http://localhost:29193/api/v1/data"
    protocol   = "http://localhost:29193/api/v1/ids"
    public     = "http://localhost:29193/public"
    control    = "http://localhost:29193/control"
  }
}

resource "edc_asset" "s3" {
  asset = {
    "asset:prop:name" : "S3 with generated id",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    s3 = {
      name              = "test"
      bucket_name       = "test"
      access_key_id     = "dummy_key"
      secret_access_key = "dummy_key"
    }
  }
}


resource "edc_policy" "policy" {
  id = "newPolicy"
  policy = {
    uid = "231802-bb34-11ec-8422-0242ac120002",
    permissions = [
      {
        edctype = "dataspaceconnector:permission",
        target  = "assetId",
        action = {
          type = "USE"
        },
      }
    ]
  }
}


resource "edc_contract_definition" "contract" {
  access_policy_id   = edc_policy.policy.id
  contract_policy_id = edc_policy.policy.id
  validity           = 600
  criteria = [
    {
      operand_left  = "asset:prop:id"
      operator      = "="
      operand_right = edc_asset.s3.id
    }
  ]
}

data "edc_contract_definition" "my_contract" {
  id = edc_contract_definition.contract.id
}

output "my_contract_id" {
  value = data.edc_contract_definition.my_contract.id
}
