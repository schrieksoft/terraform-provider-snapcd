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

data "snapcd_module" "mymodule" {
  name         = "mymodule"
  namespace_id = data.snapcd_namespace.mynamespace.id
}

resource "snapcd_module_mission" "propose_fix" {
  agent_id     = data.snapcd_agent.myagent.id
  module_id    = data.snapcd_module.mymodule.id
  name         = "propose-fix-mymodule"
  mission_type = "ProposeFix"
}
