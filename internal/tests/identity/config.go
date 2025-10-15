package identity

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var GroupMemberCreateConfig = `
resource "snapcd_group_member" "this" { 
  group_id  	 		     = snapcd_group.this.id
  principal_id   		     = data.snapcd_service_principal.this.id
  group_member_discriminator = "ServicePrincipal"
}`

var GroupCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" { 
  name  = "somevalue%s"
}`)

var AnotherGroupCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "another" { 
  name  = "anothervalue%s"
}`)

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debug"
}
`

var UserGetConfig = `
data "snapcd_user" "this" { 
  user_name  = "kschriek@gmail.com"
}`
