terraform {
  required_providers {
    edc = {
      source  = "Think-iT-Labs/edc"
      version = "0.1.0"
    }
  }
}

provider "edc" {
  alias = "org_a"
  token = "1234"
  addresses = {
    default    = "http://localhost:29193/api"
    management = "http://localhost:29193/api/v1/data"
    protocol   = "http://localhost:29193/api/v1/ids"
    public     = "http://localhost:29193/public"
    control    = "http://localhost:29193/control"
  }
}

provider "edc" {
  alias = "org_b"
  token = "1234"
  addresses = {
    default    = "http://localhost:28183/api"
    management = "http://localhost:28183/api/v1/data"
    protocol   = "http://localhost:28183/api/v1/ids"
    public     = "http://localhost:28183/public"
    control    = "http://localhost:28183/control"
  }
}

provider "edc" {
  alias = "org_c"
  token = "1234"
  addresses = {
    default    = "http://localhost:27173/api"
    management = "http://localhost:27173/api/v1/data"
    protocol   = "http://localhost:27173/api/v1/ids"
    public     = "http://localhost:27173/public"
    control    = "http://localhost:27173/control"
  }
}
