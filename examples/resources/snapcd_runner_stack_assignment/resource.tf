data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

data "snapcd_service_principal" "mysp" {
  client_id = "mysp"
}

resource "snapcd_runner_service_principal_assignment" "myrunner_mysp" {
  runner_id = data.snapcd_runner.myrunner.id
  stack_id  = data.snapcd_service_principal.mysp.id
}
