data "snapcd_agent" "myagent" {
  name = "myagent"
}

data "snapcd_stack" "mystack" {
  name = "mystack"
}

resource "snapcd_stack_mission" "recommend_approvals" {
  agent_id     = data.snapcd_agent.myagent.id
  stack_id     = data.snapcd_stack.mystack.id
  name         = "recommend-approvals-mystack"
  mission_type = "ApprovalRecommend"
}
