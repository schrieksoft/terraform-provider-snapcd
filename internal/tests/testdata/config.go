package testdata

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var StackCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack" "this" {
  name  = "somevalue%s"
}`)

var NamespaceCreateConfig = StackCreateConfig + providerconfig.AppendRandomString(`

resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = snapcd_stack.this.id
  default_apply_approval_threshold = 1
}
`)

var ModuleCreateConfigDelta = providerconfig.AppendRandomString(`

data "snapcd_runner" "debug" {
  name = "debug"
}

resource "snapcd_module" "this" {
  name                         	 = "somevalue%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_id                 = data.snapcd_runner.debug.id
  source_subdirectory  	         = "modules/module1"
  source_url                     = "foo"
  source_revision                = "main"
  trigger_on_definition_changed          = false
  trigger_on_upstream_output_changed     = false
  trigger_on_source_changed              = false
  trigger_on_source_changed_notification = false
  apply_approval_threshold               = 1
}
`)

var ModuleCreateConfig = NamespaceCreateConfig + ModuleCreateConfigDelta

var ModuleCreateConfigDeltaTwo = providerconfig.AppendRandomString(`

resource "snapcd_module" "two" {
  name                         	 = "somevalueTwo%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_id                 = data.snapcd_runner.debug.id
  source_subdirectory  	         = "modules/module1"
  source_url                     = "foo"
  source_revision                = "main"
  trigger_on_definition_changed          = false
  trigger_on_upstream_output_changed     = false
  trigger_on_source_changed              = false
  trigger_on_source_changed_notification = false
  apply_approval_threshold               = 1
}
`)

var RunnerCreateConfig = providerconfig.AppendRandomString(`
data "snapcd_service_principal" "runner" {
	client_id  = "debugRunner"
}

resource "snapcd_runner" "this" {
  name  			   = "somevalue%s"
  service_principal_id = data.snapcd_service_principal.runner.id
}`)

var AgentServicePrincipalConfig = `
data "snapcd_service_principal" "agent" {
	client_id  = "debug"
}
`

var AgentCreateConfig = AgentServicePrincipalConfig + providerconfig.AppendRandomString(`
resource "snapcd_agent" "this" {
  name                 = "somevalue%s"
  service_principal_id = data.snapcd_service_principal.agent.id
}`)

var StackSecretDebugDataSourceDelta = `
data "snapcd_stack_secret" "debug" {
	name 	  = "debug"
    stack_id = data.snapcd_stack.debug.id
}
`

var StackDebugDataSourceDelta = `
data "snapcd_stack" "debug" {
  name = "debug"
}
`
var NamespaceDebugDataSourceDelta = `
data "snapcd_namespace" "debug" {
  name = "debug"
  stack_id = data.snapcd_stack.debug.id
}
`

var ModuleDebugDataSourceDelta = `
data "snapcd_module" "debug" {
  name = "debug"
  namespace_id = data.snapcd_namespace.debug.id
}
`
