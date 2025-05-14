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
	DescSharedEngine                = "Determines which binary will be used during deployment. Must be one of 'OpenTofu' and 'Terraform'. Setting this to 'OpenTofu' will use `tofu`. Setting it to 'Terraform' will use `terraform`. "
	DescSharedOutputSecretStoreId   = "The ID of the Secret Store that will be used to store this Module's outputs. Note that for an 'Output' step to successfully use this Secret Store, it must either be deployed as `is_assigned_to_all_scopes=true`, or assigned via module/namespace/stack assignment. "
)
