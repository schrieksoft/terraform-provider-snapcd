// SPDX-License-Identifier: MPL-2.0

package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceAwsSecretStore_Create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AwsSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_aws_secret_store.this", "id"),
				),
			},
		},
	})
}

func TestAccResourceAwsSecretStore_CreateUpdate(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AwsSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_aws_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_aws_secret_store.this", "name", providerconfig.AppendRandomString("somevalue%s")),
				),
			},
			{
				Config: providerconfig.ProviderConfig + providerconfig.AppendRandomString(`
resource "snapcd_aws_secret_store" "this" { 
  name = "someNEWvalue%s"
  region = "us-east-1"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_aws_secret_store.this", "id"),
					resource.TestCheckResourceAttr("snapcd_aws_secret_store.this", "name", providerconfig.AppendRandomString("someNEWvalue%s")),
				),
			},
		},
	})
}

func TestAccResourceAwsSecretStore_Import(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerconfig.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerconfig.ProviderConfig + AwsSecretStoreCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snapcd_aws_secret_store.this", "id"),
				),
			},
			{
				ResourceName:      "snapcd_aws_secret_store.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
