package namespace

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
)

var StackCreateConfig = testdata.StackCreateConfig

var NamespaceCreateConfig = testdata.NamespaceCreateConfig

var NamespaceUpdateConfig = StackCreateConfig + providerconfig.AppendRandomString(`
resource "snapcd_namespace" "this" {
  name                      = "somevalue%s"
  stack_id			     		    = snapcd_stack.this.id
  default_apply_approval_threshold = 2
}

`)
