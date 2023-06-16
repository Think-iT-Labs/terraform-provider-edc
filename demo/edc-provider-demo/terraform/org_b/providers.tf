terraform {
  required_providers {
    edc = {
      source  = "Think-iT-Labs/edc"
      version = "0.1.0"
    }
  }
}

provider "edc" {
  token = "1234"
  addresses = {
    default    = "http://localhost:28183/api"
    management = "http://localhost:28183/api/v1/data"
    protocol   = "http://localhost:28183/api/v1/ids"
    public     = "http://localhost:28183/public"
    control    = "http://localhost:28183/control"
  }
}
