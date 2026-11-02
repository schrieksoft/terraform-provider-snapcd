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

  # Nothing downstream reads these Outputs unless they are supplied explicitly.
  output_set_downstream_restriction = "All"
}

resource "snapcd_module" "consumer" {
  name                = "consumer"
  namespace_id        = snapcd_namespace.consumers.id
  source_revision     = "main"
  source_url          = "https://github.com/schrieksoft/snapcd-samples.git"
  source_subdirectory = "getting-started/two-module-dag/module2"
  runner_id           = data.snapcd_runner.myrunner.id
}

resource "snapcd_output_set_module_supply" "producer_to_consumer" {
  module_id           = snapcd_module.producer.id
  consuming_module_id = snapcd_module.consumer.id
}
