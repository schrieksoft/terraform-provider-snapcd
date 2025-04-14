package module_param

const (
	DescSharedId             = "Unique ID of the "
	DescSharedName1          = "Name of the "
	DescSharedName2          = " Must be unique in combination with `module_id`."
	DescSharedModuleId1      = "ID of the "
	DescSharedModuleId2      = "'s parent Module."
	DescSharedDefinitionName = "Name of the Definition from which to get take the input. Can be one of 'ModuleId', 'NamespaceId', 'StackId', 'ModuleName', 'NamespaceName', 'StackName', 'SourceUrl', 'SourceRevision' or 'SourceSubdirectory'."
	DescSharedLiteralType    = "Type of literal input."
	DescSharedLiteralValue   = "Literal value of the input."
	DescSharedModuleName     = "Name of the parent Module from which to source the take the Output."
	DescSharedNamespaceName  = "Name of the parent Namespace of the parent Module from which to take the Output."
	DescSharedReferenceName  = "Name of the Namespace Input to pull in to take as input."
	DescSharedOutputName     = "Name of Output to take as input."
	DescSharedSecretName     = "Name of the Secret to take as input."
	DescSharedSecretScope    = "Scope on which the Secret to take as input as been stored (Stack, Namespace or Module)."
	DescSharedSecretType     = "Type of literal input the secret value should be formatted as."
)
