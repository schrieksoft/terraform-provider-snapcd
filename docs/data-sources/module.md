---
page_title: "snapcd_module Data Source - snapcd"
subcategory: "Modules"
description: |-
  Use this data source to access information about an existing Module in Snap CD.
---

# snapcd_module (Data Source)

Use this data source to access information about an existing Module in Snap CD.


## Example Usage

```terraform
data "snapcd_stack" "default" {
  name = "default"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.default.id
}
data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Module. Must be unique in combination with `namespace_id`.
- `namespace_id` (String) ID of the Module's parent Namespace.

### Read-Only

- `apply_after_hook` (String) Shell script that should be executed after the 'Apply' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `apply_before_hook` (String) Shell script that should be executed before the 'Apply' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `depends_on_modules` (List of String) A list on Snap CD Modules that this Module depends on. Note that Snap CD will automatically discover depedencies based on the Module using as inputs the outputs from another Module, so use `depends_on_modules` where you want to explicitly establish a dependency where outputs are not referenced as inputs.
- `destroy_after_hook` (String) Shell script that should be executed after the 'Destroy' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `destroy_before_hook` (String) Shell script that should be executed before the 'Destroy' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `engine` (String) Determines which binary will be used during deployment. Must be one of 'OpenTofu' and 'Terraform'. Setting this to 'OpenTofu' will use `tofu`. Setting it to 'Terraform' will use `terraform`. Setting this will override any default value set on the Module's parent Namespace.
- `id` (String) Unique ID of the Module.
- `init_after_hook` (String) Shell script that should be executed after the 'Init' step of any deployment is run.Setting this will override any default value set on the Module's parent Namespace.
- `init_backend_args` (String) Arguments to pass to the 'init' command in order to set the backend. This should be a text block.Setting this will override any default value set on the Module's parent Namespace.
- `init_before_hook` (String) Shell script that should be executed before the 'Init' step of any deployment is run.Setting this will override any default value set on the Module's parent Namespace.
- `output_after_hook` (String) Shell script that should be executed after the 'Output' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `output_before_hook` (String) Shell script that should be executed before the 'Output' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `output_secret_store_id` (String) The ID of the Secret Store that will be used to store this Module's outputs. Note that for an 'Output' step to successfully use this Secret Store, it must either be deployed as `is_assigned_to_all_scopes=true`, or assigned via module/namespace/stack assignment. Setting this will override any default value set on the Module's parent Namespace.
- `plan_after_hook` (String) Shell script that should be executed after the 'Plan' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `plan_before_hook` (String) Shell script that should be executed before the 'Plan' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `plan_destroy_after_hook` (String) Shell script that should be executed after the 'PlanDestroy' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `plan_destroy_before_hook` (String) Shell script that should be executed before the 'PlanDestroy' step of any deployment is run. Setting this will override any default value set on the Module's parent Namespace.
- `runner_pool_id` (String) ID of the Runner Pool that will receive the instructions when triggering a deployment on this Module.
- `runner_self_declared_name` (String) Name of the Runner to select (should unique identify the Runner within the Runner Pool). If null a random Runner will be selected from the Runner pool on every deployment.
- `source_revision` (String) Remote revision (e.g. version number, branch, commit or tag) where the source module code is found.
- `source_revision_type` (String) How Snap CD should interpret the `source_revision` field. Must be one of 'Default' or 'SemanticVersionRange'. Setting to 'Default' means Snap CD will interpret the revision type based on the source type (for example, for a 'Git' source type it will automatically figure out whether the `source_revision` refers to a branch, tag or commit). Setting to 'SemanticVersionRange' means that Snap CD will resolve the revision to a semantic version line `vX.Y.Z` (alternatively witout the 'v' prefix of that is how your semantic version are tagged, i.e. 'X.Y.Z'). It will take the highest version within the major or minor version range that you specify. For example, specify `v2.20.*` or `v2.*`. You can also specify a specific semantic version here, e.g. `v2.20.7`. In that case the behaviour is the same as with when using 'Default', except that only valid semantic versions are accepted. NOTE that 'SemanticVersionRange' is currently only supported in combination with the 'Git' `source_type`.
- `source_subdirectory` (String) Subdirectory where the source module code is found.
- `source_type` (String) The type of remote module store that the source module code should be retrieved from. Must be one of 'Git' or 'Registry'
- `source_url` (String) Remote URL where the source module code is found.
- `trigger_on_definition_changed` (Boolean) Defaults to 'true'. If 'true', the Module will automatically be applied if its definition changes. A definition change results from fields on the Module itself, on any of its Inputs (Param or Env Var) or Extra Files being altered. So too changes to its Namespace (including Inputs and Extra Files) or Stack. Note however that Namespace and Stack changes are not notified by default. This behaviour can be changed in `snapcd_namespace` and `snapcd_stack` resource definitions.
- `trigger_on_source_changed` (Boolean) Defaults to 'true'. If 'true', the Module will automatically be applied if the source it is referencing has changed. For example, if tracking a Git branch: a new commit would constitute a change.
- `trigger_on_source_changed_notification` (Boolean) Defaults to 'false'. If 'true', the Module will automatically be applied if the 'api/Hooks/SourceChanged' endpoint is called for this Module. Use this if you want to use external tooling to inform Snap CD that a source has been changed. Consider setting `trigger_on_definition_changed` to 'false' when setting `trigger_on_definition_changed_hook` to 'true'
- `trigger_on_upstream_output_changed` (Boolean) Defaults to 'true'. If 'true', the Module will automatically be applied if Outputs from other Modules that it is referencing as Inputs (Param or Env Var) has changed.
