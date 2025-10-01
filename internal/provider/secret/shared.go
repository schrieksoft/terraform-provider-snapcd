package secret

const (
	DescId          = "Unique ID of the Secret."
	DescName        = "Unique Name within of the Secret within the Secret Store."
	DescModuleId    = "Id of the Module to scope the Secret to"
	DescNamespaceId = "Id of the Namespace to scope the Secret to"
	DescStackId     = "Id of the Stack to scope the Secret to"
	DescValue       = "Value of the to store in the Simple Secret Store. NOTE that this value **will** end up in the .tfstate file. If you wish to avoid this, create the secret directly via the API or Dashboard instead."
)
