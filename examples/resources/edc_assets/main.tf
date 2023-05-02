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

resource "edc_asset" "custom" {
  asset = {
    "asset:prop:id" : "customAssetId",
    "asset:prop:name" : "customAssetName",
    "asset:prop:contenttype" : "application/json",
  }

  data = {
    custom = <<EOF
    {
      "type"              : "amazonS3",
      "name"              : "testCustom",
      "bucket_name"       : "testCustom",
      "access_key_id"     : "dummy_key_custom",
      "secret_access_key" : "dummy_key_custom"
    }
    EOF
  }
}
