# Copyright (c) HashiCorp, Inc.

provider "snapcd" {
  url                           = "https://localhost:20002"
  insecure_skip_verify          = true
  health_check_interval_seconds = 5
  health_check_timeout_seconds  = 20
  client_id                     = "Terraformer"
  client_secret                 = "somesecretforsampleclient"
}