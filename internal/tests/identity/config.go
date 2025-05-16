package identity

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var GroupMemberCreateConfig = `
resource "snapcd_group_member" "this" { 
  group_id  	 		  = snapcd_group.this.id
  principal_id   		  = data.snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
}`

var GroupCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" { 
  name  = "somevalue%s"
}`)

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "IntegratedRunner"
}
`

var UserGetConfig = `
data "snapcd_user" "this" { 
  user_name  = "kschriek@gmail.com"
}`
