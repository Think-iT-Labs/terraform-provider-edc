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
