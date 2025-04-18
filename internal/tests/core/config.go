package core

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var ModuleCreateConfig = NamespaceCreateConfig + providerconfig.AppendRandomString(`

data "snapcd_runner_pool" "default" {
  name = "default"
}

resource "snapcd_module" "this" {
  name                         	 = "somevalue%s"
  namespace_id                	 = snapcd_namespace.this.id
  runner_pool_id                 = data.snapcd_runner_pool.default.id
  source_subdirectory  	       = "modules/module1"
  source_url                     = "foo"
  init_before_hook				       = "fooBeforeHook"
}
`)

var NamespaceCreateConfig = providerconfig.AppendRandomString(`
data "snapcd_stack" "default" {
  name  = "default"
}

resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = data.snapcd_stack.default.id
  default_init_before_hook  = "foo"
}
`)

var NamespaceUpdateConfig = providerconfig.AppendRandomString(`


data "snapcd_stack" "default" {
  name  = "default"
}

resource "snapcd_namespace" "this" {
  name                        	 = "somevalue%s"
  stack_id			     		 = data.snapcd_stack.default.id
  default_init_before_hook       = "bar"
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
