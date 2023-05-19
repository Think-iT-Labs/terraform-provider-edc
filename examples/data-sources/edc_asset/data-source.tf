terraform {
  required_providers {
    edc = {
      source = "think-it-labs/edc"
    }
  }
}

provider "edc" {
  token = "test-token"
  addresses = {
    default    = "http://localhost:29193/api"
    management = "http://localhost:29193/api/v1/data"
    protocol   = "http://localhost:29193/api/v1/ids"
    public     = "http://localhost:29193/public"
    control    = "http://localhost:29193/control"
  }
}

data "edc_asset" "asset" {
  id = "6dca15a6-8cb4-4f63-8739-e6adce7430ad"
}


output "asset_output" {
  value = data.edc_asset.asset
}
