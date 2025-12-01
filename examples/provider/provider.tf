# Authenticate with Service Princpal

variable "client_id" {}
variable "client_secret" { sensitive = true }
variable "organization_id" {}

provider "snapcd" {
  client_id       = var.client_id
  client_secret   = var.client_secret
  organization_id = var.organization_id
}

# Authenticate with Access Token

variable "access_token" { sensitive = true }
variable "organization_id" {}

provider "snapcd" {
  access_token    = var.access_token
  organization_id = var.organization_id
}
