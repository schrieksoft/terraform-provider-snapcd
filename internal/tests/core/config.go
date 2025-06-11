package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var ModuleCreateConfigDelta = providerconfig.AppendRandomString(`

data "snapcd_runner_pool" "default" {
  name = "default"
}

resource "snapcd_module" "this" {
  name                         	 = "somevalue%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_pool_id                 = data.snapcd_runner_pool.default.id
  source_subdirectory  	         = "modules/module1"
  source_url                     = "foo"
  source_revision                = "main"
  init_before_hook				       = "fooBeforeHook"
  trigger_on_definition_changed          = false
  trigger_on_upstream_output_changed     = false
  trigger_on_source_changed              = false
  trigger_on_source_changed_notification = false
  apply_approval_threshold               = 1
}
`)

var ModuleCreateConfigDeltaTwo = providerconfig.AppendRandomString(`

resource "snapcd_module" "two" {
  name                         	 = "somevalueTwo%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_pool_id                 = data.snapcd_runner_pool.default.id
  source_subdirectory  	         = "modules/module1"
  source_url                     = "foo"
  source_revision                = "main"
  init_before_hook				       = "fooBeforeHook"
  trigger_on_definition_changed          = false
  trigger_on_upstream_output_changed     = false
  trigger_on_source_changed              = false
  trigger_on_source_changed_notification = false
  apply_approval_threshold               = 1
}
`)

var ModuleCreateConfig = NamespaceCreateConfig + ModuleCreateConfigDelta

var NamespaceCreateConfig = providerconfig.AppendRandomString(`
data "snapcd_stack" "default" {
  name  = "default"
}

resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = data.snapcd_stack.default.id
  default_init_before_hook  = "foo"
  default_apply_approval_threshold = 1
}
`)

var NamespaceUpdateConfig = providerconfig.AppendRandomString(`


data "snapcd_stack" "default" {
  name  = "default"
}

resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = data.snapcd_stack.default.id
  default_init_before_hook  = "bar"
}

`)

var RunnerPoolCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_runner_pool" "this" {
  name  = "somevalue%s"
}`)

var StackCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack" "this" {
  name  = "somevalue%s"
}`)

var PrexistingStack = `
resource "snapcd_stack" "this" {
  name  = "default"
}`

var SourceRefresherPreselectionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_source_refresher_preselection" "this" {
  source_url     = "somevalue%s"
  runner_pool_id = snapcd_runner_pool.this.id
}`)
