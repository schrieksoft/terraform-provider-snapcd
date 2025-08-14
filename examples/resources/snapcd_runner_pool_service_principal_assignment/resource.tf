data "snapcd_service_principal" "mysp" {
  client_id = "mysp"
}

resource "snapcd_runner_pool" "myrunnerpool" {
  name = "myrunnerpool"
}

resource "snapcd_runner_pool_service_principal_assignment" "myrunnerpool_mysp" {
  runner_pool_id       = snapcd_runner_pool.myrunnerpool.id
  service_principal_id = data.snapcd_service_principal.mysp.id
}
