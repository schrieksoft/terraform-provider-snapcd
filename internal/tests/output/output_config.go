package output

import "terraform-provider-snapcd/internal/tests/providerconfig"

// A consuming Namespace and Module beside the ones ModuleCreateConfig already builds, so a supply
// has somewhere to point that is not its own producer.
var ConsumerCreateConfig = providerconfig.AppendRandomString(`

resource "snapcd_namespace" "consumer" {
  name                             = "consumer%[1]s"
  stack_id                         = snapcd_stack.this.id
  default_apply_approval_threshold = 1
}

resource "snapcd_module" "consumer" {
  name                                   = "consumer%[1]s"
  namespace_id                           = snapcd_namespace.consumer.id
  runner_id                              = data.snapcd_runner.debug.id
  source_subdirectory                    = "modules/module2"
  source_url                             = "foo"
  source_revision                        = "main"
  trigger_on_definition_changed          = false
  trigger_on_upstream_output_changed     = false
  trigger_on_source_changed              = false
  trigger_on_source_changed_notification = false
  apply_approval_threshold               = 1
  state_migration_approval_threshold     = 1
}
`)

var ModuleOutputDownstreamRestrictionCreateConfig = `
resource "snapcd_module_output_downstream_restriction" "this" {
  module_id   = snapcd_module.this.id
  output_name = "db_password"
}`

var OutputSetModuleSupplyCreateConfig = `
resource "snapcd_output_set_module_supply" "this" {
  module_id           = snapcd_module.this.id
  consuming_module_id = snapcd_module.consumer.id
}`

var OutputSetNamespaceSupplyCreateConfig = `
resource "snapcd_output_set_namespace_supply" "this" {
  module_id              = snapcd_module.this.id
  consuming_namespace_id = snapcd_namespace.consumer.id
}`

var OutputModuleSupplyCreateConfig = `
resource "snapcd_output_module_supply" "this" {
  module_id           = snapcd_module.this.id
  output_name         = "vpc_id"
  consuming_module_id = snapcd_module.consumer.id
}`

var OutputNamespaceSupplyCreateConfig = `
resource "snapcd_output_namespace_supply" "this" {
  module_id              = snapcd_module.this.id
  output_name            = "vpc_id"
  consuming_namespace_id = snapcd_namespace.consumer.id
}`
