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
    default    = "http://localhost:27173/api"
    management = "http://localhost:27173/api/v1/data"
    protocol   = "http://localhost:27173/api/v1/ids"
    public     = "http://localhost:27173/public"
    control    = "http://localhost:27173/control"
  }
}
