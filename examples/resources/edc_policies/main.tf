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

resource "edc_policy" "pol" {
  id = "abcPolicy"
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
