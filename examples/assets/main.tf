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
    "asset:prop:id" : "assetId",
    "asset:prop:name" : "assetName",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    s3 = {
      type              = "AmazonS3"
      name              = "test"
      bucket_name       = "test"
      access_key_id     = "dummy_key"
      secret_access_key = "dummy_key"
    }
  }
}
