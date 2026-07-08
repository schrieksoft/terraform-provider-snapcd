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

var ServicePrincipalDataSourceConfig = `
data "snapcd_service_principal" "this" {
	client_id = "debugTestTarget1"
}
`

var UserGetConfig = `
data "snapcd_user" "this" {
  user_name  = "debug@preseeded.io"
}`
