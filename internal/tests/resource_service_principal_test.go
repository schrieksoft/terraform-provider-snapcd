// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var servicePrincipalCreateConfig = appendRandomString(`
resource "snapcd_service_principal" "this" { 
  client_id  	 = "somevalue%s"
  client_secret  = "veryverysecret"
}`)


func TestAccResourceServicePrincipal_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_Create_CreateAgain(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
			{
				Config: providerConfig + servicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
resource "snapcd_service_principal" "this" { 
  client_id  	 = "someNEWvalue%s"
  client_secret  = "veryverysecret"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", appendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceServicePrincipal_CreateUpdateNewSecret(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + servicePrincipalCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_service_principal.this", "id"),
					resource.TestCheckResourceAttr("snapcd_service_principal.this", "client_id", appendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerConfig + appendRandomString(`
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
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: providerConfig + servicePrincipalCreateConfig,
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
