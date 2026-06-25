package stack

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var StackCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_stack" "this" {
  name  = "somevalue%s"
}`)

var PrexistingStack = `
resource "snapcd_stack" "should_fail" {
  name  = "debug"
}`
