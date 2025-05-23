---
page_title: "snapcd_group Data Source - snapcd"
subcategory: "Identity Access Management"
description: |-
  Use this data source to access information about an existing Group in Snap CD.
---

# snapcd_group (Data Source)

Use this data source to access information about an existing Group in Snap CD.


## Example Usage

```terraform
data "snapcd_group" "mygroup" {
  name = "mygroup"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Unique Name of the Group.

### Read-Only

- `description` (String) Description of the Group.
- `id` (String) Unique ID of the Group.
