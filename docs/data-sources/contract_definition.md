---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "edc_contract_definition Data Source - terraform-provider-edc"
subcategory: ""
description: |-
  Contract Definition Data Source
---

# edc_contract_definition (Data Source)

Contract Definition Data Source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Contract definition identifier

### Read-Only

- `access_policy_id` (String) Access policy identifier
- `contract_policy_id` (String) Contract policy identifier
- `created_at` (Number) Created at timestamp in seconds
- `criteria` (Attributes List) (see [below for nested schema](#nestedatt--criteria))
- `validity` (Number) Validity

<a id="nestedatt--criteria"></a>
### Nested Schema for `criteria`

Read-Only:

- `operand_left` (String)
- `operand_right` (String)
- `operator` (String)
