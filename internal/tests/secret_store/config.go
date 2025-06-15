package secret_store

import (
	"terraform-provider-snapcd/internal/tests/providerconfig"
)

var AzureKeyVaultSecretStoreCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_azure_key_vault_secret_store" "this" { 
  name  = "somevalue%s"
  key_vault_url = "https://snapcdlocaltesting.vault.azure.net/"
}`)

var SimpleSecretStoreCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_simple_secret_store" "this" { 
  name  = "somevalue%s"
}`)

var AwsSecretsManagerSecretStoreCreateConfig = providerconfig.AppendRandomString(`
resource "snapcd_aws_secrets_manager_secret_store" "this" { 
  name  = "somevalue%s"
  region = "us-east-1"
}`)
