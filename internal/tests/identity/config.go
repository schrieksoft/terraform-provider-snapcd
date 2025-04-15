package identity

import (
	providerconfig "terraform-provider-snapcd/internal/tests/providerconfig"
)

var GroupMemberCreateConfig = `
resource "snapcd_group_member" "this" { 
  group_id  	 		  = snapcd_group.this.id
  principal_id   		  = snapcd_service_principal.this.id
  principal_discriminator = "ServicePrincipal"
}`

var ServicePrincipalCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" { 
  client_id  	 = "somevalue%s"
  client_secret  = "veryverysecret"
}`)

var GroupCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" { 
  name  = "somevalue%s"
}`)

var UserGetConfig = `
data "snapcd_user" "this" { 
  user_name  = "kschriek@gmail.com"
}`
