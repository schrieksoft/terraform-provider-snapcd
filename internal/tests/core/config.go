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
// target_repo_revision         	 = "main"
// target_repo_url              	 = "git@github.com:karlschriek/tf-samples.git"
  target_module_relative_path  	 = "modules/module1"
  provider_cache_enabled         = true
  module_cache_enabled         	 = true
//  depends_on_modules		 	 = []
  select_on           			 = "PoolId"
  select_strategy     			 = "FirstOf"
  init_before_hook				 = "fooBeforeHook"
  //selected_consumer_id 		 = "cli"
}
`)

var NamespaceCreateConfig = providerconfig.AppendRandomString(`
data "snapcd_stack" "default" { 
  name  = "default"
}

resource "snapcd_namespace" "this" { 
  name                        	 = "somevalue%s"
  stack_id			     		 = data.snapcd_stack.default.id
  default_init_before_hook       = "foo"
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
