package core

import (
	"terraform-provider-snapcd/internal/tests/identity"
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var DependsOnModuleCreateConfig = `
resource "snapcd_depends_on_module" "this" { 
  module_id = snapcd_module.this.id
  depends_on_module_id = snapcd_module.two.id
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

var StackSecretDebugDataSourceDelta = `
data "snapcd_stack_secret" "debug" {
	name 	  = "debug"
    stack_id = data.snapcd_stack.debug.id
}
`

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
  runner_id                 = data.snapcd_runner.debug.id
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

var ModuleCreateConfigDeltaThree = providerconfig.AppendRandomString(`

resource "snapcd_module" "three" {
  name                         	 = "somevalueThree%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_id                 = data.snapcd_runner.debug.id
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

var NamespaceCreateConfig = StackCreateConfig + providerconfig.AppendRandomString(`

resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = snapcd_stack.this.id
  default_init_before_hook  = "foo"
  default_apply_approval_threshold = 1
}
`)

var NamespaceUpdateConfig = StackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = snapcd_stack.this.id
  default_init_before_hook  = "bar"
}

`)

var RunnerCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_runner" "this" {
  name  = "somevalue%s"
}`)

var RunnerCreateConfigWithThreshold = providerconfig.AppendRandomString(`
resource "snapcd_runner" "this" {
  name  = "somevalue%s"
  custom_command_approval_threshold = 2
}`)

var StackCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack" "this" {
  name  = "somevalue%s"
}`)

var PrexistingStack = `
data "snapcd_stack" "debug" {
  name  = "debug"
}`

var SourceRefresherPreselectionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_source_refresher_preselection" "this" {
  source_url     = "somevalue%s"
  runner_id = snapcd_runner.this.id
}`)

var CustomCommandPreApprovalCreateConfig = RunnerCreateConfig + identity.ServicePrincipalDataSourceConfig + `

resource "snapcd_custom_command_pre_approval" "this" {
  runner_id                   = snapcd_runner.this.id
  command_text                     = "terraform plan"
  approver_principal_id            = data.snapcd_service_principal.this.id
  approver_principal_discriminator = "ServicePrincipal"
}`
