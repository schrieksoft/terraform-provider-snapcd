resource "snapcd_runner" "myrunnerpool" {
  name = "myrunnerpool"
}

data "snapcd_service_principal" "mysp" {
  client_id = "mysp"
}

resource "snapcd_runner_service_principal_assignment" "myrunnerpool_mysp" {
  runner_id = snapcd_runner.myrunnerpool.id
  stack_id  = data.snapcd_service_principal.mysp.id
}
