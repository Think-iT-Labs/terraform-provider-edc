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
    default    = "http://localhost:19191/api"
    control    = "http://localhost:19192/control"
    management = "http://localhost:19193/api/v1/data"
    protocol   = "http://localhost:19194/api/v1/ids"
    public     = "http://localhost:19291/public"
  }
}
