package state_store

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var StateStoreCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_state_store" "this" {
  name = "somevalue%s"
}`)
