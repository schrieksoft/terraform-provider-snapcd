// SPDX-License-Identifier: MPL-2.0

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
  name  = "grp_create_%s"
}`)

var GroupImportConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_import_%s"
}`)

var GroupUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_update_%s"
}`)

var GroupUpdateNewConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_update_new_%s"
}`)

var GroupDatasourceConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_ds_%s"
}`)

var GroupMemberCreateGroupConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_member_create_%s"
}`)

var GroupMemberUpdateGroupConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_member_update_%s"
}`)

var GroupMemberImportGroupConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "this" {
  name  = "grp_member_import_%s"
}`)

var AnotherGroupCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_group" "another" {
  name  = "grp_member_update_another_%s"
}`)

var ServicePrincipalCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id = "sp_create_%s"
}`)

var ServicePrincipalUpdateConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id   = "sp_update_%s"
  is_disabled = true
}`)

var ServicePrincipalImportConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id = "sp_import_%s"
}`)

var ServicePrincipalWithSecretConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" {
  client_id     = "sp_secret_%s"
  client_secret = "shouldnotbeaccepted"
}`)

var FederatedCredentialServicePrincipalConfig = providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "fedcred" {
  client_id = "sp_fedcred_%s"
}`)

var FederatedCredentialCreateConfig = FederatedCredentialServicePrincipalConfig + `
resource "snapcd_service_principal_federated_credential" "this" {
  service_principal_id = snapcd_service_principal.fedcred.id
  name                 = "fedcred_create"
  issuer               = "https://oidc.example.com/cluster"
  subject              = "system:serviceaccount:snapcd:runner"
  audience             = "snapcd"
}`

var FederatedCredentialUpdateConfig = FederatedCredentialServicePrincipalConfig + `
resource "snapcd_service_principal_federated_credential" "this" {
  service_principal_id = snapcd_service_principal.fedcred.id
  name                 = "fedcred_create"
  issuer               = "https://oidc.example.com/other-cluster"
  subject              = "system:serviceaccount:snapcd:agent"
  audience             = "snapcd"
  is_disabled          = true
}`

var FederatedCredentialReplaceConfig = FederatedCredentialServicePrincipalConfig + providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "other" {
  client_id = "sp_fedcred_other_%s"
}`) + `
resource "snapcd_service_principal_federated_credential" "this" {
  service_principal_id = snapcd_service_principal.other.id
  name                 = "fedcred_create"
  issuer               = "https://oidc.example.com/cluster"
  subject              = "system:serviceaccount:snapcd:runner"
  audience             = "snapcd"
}`

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debugTestTarget1"
}
`

var UserGetConfig = `
data "snapcd_user" "this" {
  user_name  = "debug@preseeded.io"
}`
