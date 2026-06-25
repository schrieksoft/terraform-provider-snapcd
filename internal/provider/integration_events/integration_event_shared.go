package integration_events

var integrationTriggerValues = []string{
	"JobSucceeded",
	"JobFailed",
	"JobAwaitingApproval",
	"JobApproved",
	"JobDeclined",
	"JobCancelled",
	"MissionStarted",
	"MissionMilestoneReported",
	"MissionCompleted",
	"MissionFaulted",
}

const (
	DescEventId            = "Unique ID of the Integration Event."
	DescEventIntegrationId = "ID of the target integration."
	DescEventTrigger       = "Trigger this subscription fires on. Must be one of 'JobSucceeded', 'JobFailed', 'JobAwaitingApproval', 'JobApproved', 'JobDeclined', 'JobCancelled', 'MissionStarted', 'MissionMilestoneReported', 'MissionCompleted', 'MissionFaulted'."
	DescEventTemplate      = "Optional message template ({{token}} substitution). Omit to use the default for the trigger."
	DescEventFilter        = "Optional filter expression."
	DescEventIsDisabled    = "Whether the subscription is disabled."
	DescEventStackId       = "ID of the Stack this event is scoped to."
	DescEventNamespaceId   = "ID of the Namespace this event is scoped to."
	DescEventModuleId      = "ID of the Module this event is scoped to."
)

var (
	organizationIntegrationEventEndpoint = "/OrganizationIntegrationEvent"
	stackIntegrationEventEndpoint        = "/StackIntegrationEvent"
	namespaceIntegrationEventEndpoint    = "/NamespaceIntegrationEvent"
	moduleIntegrationEventEndpoint       = "/ModuleIntegrationEvent"
)
