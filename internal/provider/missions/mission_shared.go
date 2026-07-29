package missions

// MissionType values mirror the SnapCd.Contracts MissionType enum. The catalog is
// closed; these are the only valid mission_type values.
var (
	organizationMissionEndpoint = "/OrganizationMission"
	stackMissionEndpoint        = "/StackMission"
	namespaceMissionEndpoint    = "/NamespaceMission"
	moduleMissionEndpoint       = "/ModuleMission"
)
