data "snapcd_agent" "myagent" {
  name = "myagent"
}

resource "snapcd_organization_mission" "diagnose_all" {
  agent_id     = data.snapcd_agent.myagent.id
  mission_type = "AutoDiagnose" // one of: AutoDiagnose, ApprovalRecommend, SummarizeJob
}
