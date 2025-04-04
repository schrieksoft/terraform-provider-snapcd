// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)


func TestAccResourceServicePrincipal_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_Create_CreateAgain(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" { 
  client_id  	 = "someNEWvalue%s"
  client_secret  = "veryverysecret"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_CreateUpdateNewSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + ServicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_service_principal" "this" { 
  client_id  	 = "someNEWvalue%s"
  client_secret  = "veryveryNEWsecret"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_secret", "veryveryNEWsecret"),
				),
			},
		},
	})
}

// func TestAccResourceServicePrincipal_Import(t *testing.T) {
// 	resource.UnitTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerconfig.ProviderConfig + servicePrincipalCreateConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
// 				),
// 			},
// 			{
// 				ResourceName:      "snapcd_service_principal.this",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 		},
// 	})
// }
