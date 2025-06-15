package namespace_input

const (
	DescSharedId             = "Unique ID of the "
	DescSharedName1          = "Name of the "
	DescSharedName2          = " Must be unique in combination with `namespace_id`."
	DescSharedNamespaceId1   = "ID of the "
	DescSharedNamespaceId2   = "'s parent Namespace."
	DescSharedDefinitionName = "Name of the Definition from which to get take the input. Must be one of 'ModuleId', 'NamespaceId', 'StackId', 'ModuleName', 'NamespaceName', 'StackName', 'SourceUrl', 'SourceRevision' and 'SourceSubdirectory'"
	DescSharedLiteralType    = "Type of literal input. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc."
	DescSharedLiteralValue   = "Literal value of the input."
	DescSharedSecretName     = "Name of the Secret to take as input."
	DescSharedSecretType     = "Type of literal input the secret value should be formatted as. Must be one of 'String' and 'NotString'. Use 'NotString' for values such as numbers, bools, list, maps etc."
	DescSharedSecretScope    = "Scope on which the Secret to take as input as been stored. Must be one of 'Stack' or 'Namespace'."
	DescSharedUsage          = "Whether the input should be used by default on all Modules, or only when explicitly selected on the Module itself. Must be one of 'UseIfSelected' and 'UseByDefault'"
	DescSharedInputKind		 = "The kind of input. Must be one of 'Param' or 'EnvVar'. Changing this will force the resource to be recreated."
)
