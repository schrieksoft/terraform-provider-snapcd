package shared

const (
	DescSharedEngine = "Determines which binary will be used during deployment. Must be one of 'OpenTofu', 'Terraform' or 'Pulumi'. Setting this to 'OpenTofu' will use `tofu`. Setting it to 'Terraform' will use `terraform`. Setting this to 'Pulumi' will use `pulumi`. "

	DescSharedCleanInitEnabled          = "Setting will remove all .terraform* files and folders (state files, locks, downloaded providers, downloaded modules etc.) and perform a clean init every time the Module is executed. "
	DescSharedDriftCheckEnabled         = "Setting this to true will periodically trigger an Apply job to check for drift in the deployed infrastructure. "
	DescSharedDriftCheckIntervalMinutes = "The number of minutes between drift checks. If not set, the system default (24 hours) is used. Note that irrespective of what is set here, these those will not be fired more regularly than the minimum internal as defined by your subscription tier. "

	DescSharedApplyApprovalThreshold   = "The number of Users (or Service Principals) that need to approve before an 'Apply' plan is executed. "
	DescSharedDestroyApprovalThreshold = "The number of Users (or Service Principals) that need to approve before an 'Destroy' plan is executed. "
	DescSharedApprovalTimeoutMinutes   = "The number of minutes a Job should remain in the 'WaitingForApproval' in the case of an 'Apply' or 'Destroy' plan that requires approval. After this time elapses the Job will be stopped and any queued Jobs will start. "
	DescZeroThreshold                  = " If set neither on Module nor on Namespace then a threshold of 0 is used."
	DescZeroTimeout                    = " If set neither on Module nor on Namespace the Jobs will wait for an approval decision indefinitely."
)
