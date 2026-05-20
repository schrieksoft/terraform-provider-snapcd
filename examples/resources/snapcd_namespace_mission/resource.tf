data "snapcd_agent" "myagent" {
  name = "myagent"
}

data "snapcd_stack" "mystack" {
  name = "mystack"
}

data "snapcd_namespace" "mynamespace" {
  name     = "mynamespace"
  stack_id = data.snapcd_stack.mystack.id
}

resource "snapcd_namespace_mission" "generate_docs" {
  agent_id     = data.snapcd_agent.myagent.id
  namespace_id = data.snapcd_namespace.mynamespace.id
  name         = "generate-docs-mynamespace"
  mission_type = "GenerateDocs"
}
