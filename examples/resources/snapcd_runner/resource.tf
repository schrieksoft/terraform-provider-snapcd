// Fetch the Service Principal you created in step 1.
data "snapcd_service_principal" "my_service_principal" {
  client_id = "MyServicePrincipal"
}

resource "snapcd_runner" "my_runner" {
  name                       = "MyGeneralRunner"
  service_principal_id       = data.snapcd_service_principal.my_service_principal.id
  is_supplied_to_all_modules = true // all modules from all stacks can use this Runner. If you want to restrict this, set to "false" and then use "snapcd_runner_stack_supply", "snapcd_runner_namespace_supply" or "snapcd_runner_module_supply" to assign to a narrower scope.
}
