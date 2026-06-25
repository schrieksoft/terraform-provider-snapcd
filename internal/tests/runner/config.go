package runner

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
	"terraform-provider-snapcd/internal/tests/testdata"
)

var RunnerCreateConfig = testdata.RunnerCreateConfig

var RunnerServicePrincipalConfig = `
data "snapcd_service_principal" "runner" {
	client_id  = "debugRunner"
}
`

var SourceRefresherPreselectionCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_source_refresher_preselection" "this" {
  source_url     = "somevalue%s"
  runner_id = snapcd_runner.this.id
}`)
