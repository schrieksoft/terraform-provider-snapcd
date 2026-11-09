data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_namespace" "producers" {
  name     = "producers"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace" "consumers" {
  name     = "consumers"
  stack_id = data.snapcd_stack.mystack.id
}

data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_module" "producer" {
  name                = "producer"
  namespace_id        = snapcd_namespace.producers.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module1"
  runner_id           = data.snapcd_runner.myrunner.id

  # Secret Outputs are withheld downstream; the restriction below withholds one more by name.
  output_set_downstream_restriction = "SensitiveOnly"
}

resource "snapcd_module" "consumer" {
  name                = "consumer"
  namespace_id        = snapcd_namespace.consumers.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_id           = data.snapcd_runner.myrunner.id
}

resource "snapcd_module_output_downstream_restriction" "internal_subnet_id" {
  module_id   = snapcd_module.producer.id
  output_name = "internal_subnet_id"
}
