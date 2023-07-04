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
    default    = "http://localhost:28191/api"
    control    = "http://localhost:28192/control"
    management = "http://localhost:28193/api/v1/data"
    protocol   = "http://localhost:28194/api/v1/ids"
    public     = "http://localhost:28291/public"
  }
}
