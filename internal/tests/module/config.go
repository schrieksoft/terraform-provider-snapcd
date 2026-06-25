package module

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
)

var ModuleCreateConfig = testdata.ModuleCreateConfig
var ModuleCreateConfigDeltaTwo = testdata.ModuleCreateConfigDeltaTwo

var DependsOnModuleCreateConfig = `
resource "snapcd_depends_on_module" "this" {
  module_id = snapcd_module.this.id
  depends_on_module_id = snapcd_module.two.id
}
`

var ModuleCreateConfigDeltaThree = providerconfig.AppendRandomString(`

resource "snapcd_module" "three" {
  name                         	 = "somevalueThree%s"
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
