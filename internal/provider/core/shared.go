package core

const (
	DescSharedInitBeforeHook = "Shell script that should be executed before the 'Init' step of any deployment is run."
	DescSharedInitAfterHook  = "Shell script that should be executed after the 'Init' step of any deployment is run."

	DescSharedInitBackedArgs        = "Arguments to pass to the 'init' command in order to set the backend. This should be a text block."
	DescSharedPlanBeforeHook        = "Shell script that should be executed before the 'Plan' step of any deployment is run. "
	DescSharedPlanAfterHook         = "Shell script that should be executed after the 'Plan' step of any deployment is run. "
	DescSharedPlanDestroyBeforeHook = "Shell script that should be executed before the 'PlanDestroy' step of any deployment is run. "
	DescSharedPlanDestroyAfterHook  = "Shell script that should be executed after the 'PlanDestroy' step of any deployment is run. "
	DescSharedApplyBeforeHook       = "Shell script that should be executed before the 'Apply' step of any deployment is run. "
	DescSharedApplyAfterHook        = "Shell script that should be executed after the 'Apply' step of any deployment is run. "
	DescSharedDestroyBeforeHook     = "Shell script that should be executed before the 'Destroy' step of any deployment is run. "
	DescSharedDestroyAfterHook      = "Shell script that should be executed after the 'Destroy' step of any deployment is run. "
	DescSharedOutputBeforeHook      = "Shell script that should be executed before the 'Output' step of any deployment is run. "
	DescSharedOutputAfterHook       = "Shell script that should be executed after the 'Output' step of any deployment is run. "
	DescSharedValidateBeforeHook    = "Shell script that should be executed before the 'Validate' step of any deployment is run. "
	DescSharedValidateAfterHook     = "Shell script that should be executed after the 'Validate' step of any deployment is run. "
	DescSharedEngine                = "Determines which binary will be used during deployment. Must be one of 'OpenTofu', 'Terraform' or 'Pulumi'. Setting this to 'OpenTofu' will use `tofu`. Setting it to 'Terraform' will use `terraform`. Setting this to 'Pulumi' will use `pulumi`. "
	DescSharedOutputSecretStoreId   = "The ID of the Secret Store that will be used to store this Module's outputs. Note that for an 'Output' step to successfully use this Secret Store, it must either be deployed as `is_assigned_to_all_scopes=true`, or assigned via module/namespace/stack assignment. "

	DescSharedAutoUpgradeEnabled     = "Deprecated: Use Terraform Flag resources instead. Setting this to true will add the `-upgrade` flag whenever `init` is called. "
	DescSharedAutoReconfigureEnabled = "Deprecated: Use Terraform Flag resources instead. Setting this to true will add the `-reconfigure` flag whenever `init` is called. "
	DescSharedAutoMigrateEnabled     = "Deprecated: Use Terraform Flag resources instead. Setting this to true will add the `-migrate-state` flag whenever `init` is called. "
	DescSharedCleanInitEnabled       = "Setting will remove all .terraform* files and folders (state files, locks, downloaded providers, downloaded modules etc.) and perform a clean init every time the Module is executed. "

	DescSharedApplyApprovalThreshold   = "The number of Users (or Service Principals) that need to approve before an 'Apply' plan is executed. "
	DescSharedDestroyApprovalThreshold = "The number of Users (or Service Principals) that need to approve before an 'Destroy' plan is executed. "
	DescSharedApprovalTimeoutMinutes   = "The number of minutes a Job should remain in the 'WaitingForApproval' in the case of an 'Apply' or 'Destroy' plan that requires approval. After this time elapses the Job will be stopped and any queued Jobs will start. "
	DescZeroThreshold                  = " If set neither on Module nor on Namespace then a threshold of 0 is used."
	DescZeroTimeout                    = " If set neither on Module nor on Namespace the Jobs will wait for an approval decision indefinitely."

	DescRunnerIsDisabled = "Indicates whether or not the Runner is disabled"
)
