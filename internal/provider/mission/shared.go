package mission

// MissionType values mirror the SnapCd.Contracts MissionType enum. The catalog is
// closed; these are the only valid mission_type values.
var missionTypeValues = []string{
	"AutoDiagnose",
	"ApprovalRecommend",
	"SummarizeJob",
}

const (
	DescMissionId          = "Unique ID of the Mission."
	DescMissionAgentId     = "ID of the Agent that runs this Mission."
	DescMissionType        = "Which named mission definition this row references. Must be one of 'AutoDiagnose', 'ApprovalRecommend' and 'SummarizeJob'."
	DescMissionSidecarName = "Optional named-sidecar override sent to the agent at dispatch. When unset (null), the agent invokes its only registered sidecar; the run fails if the agent has zero or multiple sidecars and no name was supplied."
	DescMissionIsDisabled  = "Indicates whether or not the Mission is disabled."
	DescMissionStackId     = "ID of the Stack this Mission is scoped to."
	DescMissionNamespaceId = "ID of the Namespace this Mission is scoped to."
	DescMissionModuleId    = "ID of the Module this Mission is scoped to."
)

var (
	organizationMissionEndpoint = "/OrganizationMission"
	stackMissionEndpoint        = "/StackMission"
	namespaceMissionEndpoint    = "/NamespaceMission"
	moduleMissionEndpoint       = "/ModuleMission"
)
