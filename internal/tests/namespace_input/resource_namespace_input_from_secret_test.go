// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"terraform-provider-snapcd/internal/tests/core"
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var NamespaceInputFromSecretCreateConfig = providerconfig.AppendRandomString(`
data "snapcd_stack_secret" "this" {
	name 	  = "debug"
    stack_id  = "10000000-0000-0000-0000-000000000000"
}

resource "snapcd_namespace_input_from_secret" "this" { 
  input_kind 	= "Param"
  namespace_id  = snapcd_namespace.this.id
  name  		= "somevalue%s"
  secret_id 	= data.snapcd_stack_secret.this.id
}
`)

var NamespaceInputFromSecretCreateConfigForImport = providerconfig.AppendRandomString(`
data "snapcd_stack_secret" "this" {
	name 	  = "debug"
    stack_id  = "10000000-0000-0000-0000-000000000000"
}

resource "snapcd_namespace_input_from_secret" "this" { 
  input_kind 	= "Param"
  namespace_id  = snapcd_namespace.this.id
  name  		= "someImportvalue%s"
  secret_id 	= data.snapcd_stack_secret.this.id
}
`)

var NamespaceInputFromSecretCreateConfigNew = providerconfig.AppendRandomString(`
data "snapcd_stack_secret" "this" {
	name 	  = "debug"
    stack_id  = "10000000-0000-0000-0000-000000000000"
}

resource "snapcd_namespace_input_from_secret" "this" { 
  input_kind 	= "Param"
  namespace_id  = snapcd_namespace.this.id
  name  		= "someNEWvalue%s"
  secret_id 	= data.snapcd_stack_secret.this.id
}  
`)

func TestAccResourceNamespaceInputFromSecret_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromSecret_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_secret.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfigNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
					resource.TestCheckResourceAttr("snapcd_namespace_input_from_secret.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceNamespaceInputFromSecret_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + core.NamespaceCreateConfig + NamespaceInputFromSecretCreateConfigForImport,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_namespace_input_from_secret.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_namespace_input_from_secret.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
