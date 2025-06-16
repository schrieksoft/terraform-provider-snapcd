package module_input

const (
	DescSharedId             = "Unique ID of the "
	DescSharedName1          = "Name of the "
	DescSharedName2          = " Must be unique in combination with `module_id`."
	DescSharedModuleId1      = "ID of the "
	DescSharedModuleId2      = "'s parent Module."
	DescSharedDefinitionName = "Name of the Definition from which to get take the input. Must be one of 'ModuleId', 'NamespaceId', 'StackId', 'ModuleName', 'NamespaceName', 'StackName', 'SourceUrl', 'SourceRevision' or 'SourceSubdirectory'."
	DescSharedLiteralType    = "Type of literal input. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc."
	DescSharedLiteralValue   = "Literal value of the input."
	DescSharedOutputModuleId = "ID of the Module from which to take the Output."
	DescSharedReferenceName  = "Name of the Namespace Input to pull in to take as input."
	DescSharedOutputName     = "Name of Output to take as input."
	DescSharedSecretName     = "Name of the Secret to take as input."
	DescSharedSecretType     = "The data type of the Secret to take as input. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc."
	DescSharedSecretScope    = "Scope on which the Secret to take as input as been stored. Must be one of 'Stack', 'Namespace' and 'Module'"
	DescSharedInputKind      = "The kind of input. Must be one of 'Param' or 'EnvVar'. Changing this will force the resource to be recreated."
)
