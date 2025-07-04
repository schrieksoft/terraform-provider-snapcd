---
page_title: "snapcd_hcp_secret_store Data Source - snapcd"
subcategory: "Secret Stores"
description: |-
  Use this data source to access information about an existing HCP Secret Store in Snap CD.
---

# snapcd_hcp_secret_store (Data Source)

Use this data source to access information about an existing HCP Secret Store in Snap CD.




<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Unique Name of the Secret Store.

### Read-Only

- `app_name` (String) The HCP application name.
- `id` (String) Unique ID of the Secret Store.
- `is_assigned_to_all_scopes` (Boolean) If set to true, secrets scoped to any resource in the system (any Stack, Namespace, Module or Output) can be assigned to this Secret Store
- `organization_id` (String) The HCP organization ID.
- `project_id` (String) The HCP project ID where the secret store is located.
