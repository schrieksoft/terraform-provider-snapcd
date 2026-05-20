package mission

// MissionType values mirror the SnapCd.Contracts MissionType enum. The catalog is
// closed; these are the only valid mission_type values.
var missionTypeValues = []string{
	"AutoDiagnose",
	"ApprovalRecommend",
	"ProposeFix",
	"GenerateDocs",
	"SplitMonolithicState",
}

const (
	DescMissionId          = "Unique ID of the Mission."
	DescMissionAgentId     = "ID of the Agent that runs this Mission."
	DescMissionName        = "Label for this Mission binding."
	DescMissionType        = "Which named mission definition this row references. Must be one of 'AutoDiagnose', 'ApprovalRecommend', 'ProposeFix', 'GenerateDocs' and 'SplitMonolithicState'."
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
