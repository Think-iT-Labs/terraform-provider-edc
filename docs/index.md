---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edc Provider"
subcategory: ""
description: |-
  
---

# edc Provider



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `addresses` (Attributes) (see [below for nested schema](#nestedatt--addresses))
- `token` (String)

<a id="nestedatt--addresses"></a>
### Nested Schema for `addresses`

Optional:

- `control` (String)
- `default` (String)
- `management` (String)
- `protocol` (String)
- `public` (String)
