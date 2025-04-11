variable "client_id" {}
variable "client_secret" { sensitive = true }
variable "snapcd_server_host" {}
variable "insecure_skip_verify" { default = false }

provider "snapcd" {
  url                  = "https://${var.snapcd_server_host}"
  insecure_skip_verify = var.insecure_skip_verify
  client_id            = var.client_id
  client_secret        = var.client_secret
}
