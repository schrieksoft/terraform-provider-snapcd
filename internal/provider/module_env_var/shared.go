package module_env_var

const (
	DescSharedModuleEnvVarId             = "Unique ID of the "
	DescSharedModuleEnvVarName1          = "Name of the "
	DescSharedModuleEnvVarName2          = " Must be unique in combination with `module_id`."
	DescSharedModuleEnvVarModuleId1      = "ID of the "
	DescSharedModuleEnvVarModuleId2      = "'s parent Module."
	DescSharedModuleEnvVarDefinitionName = "Name of the Definition from which to get take the input."
	DescSharedModuleEnvVarLiteralType    = "Type of literal input."
	DescSharedModuleEnvVarLiteralValue   = "Literal value of the input."
	DescSharedModuleEnvVarModuleName     = "Name of the paraent Module from which to source the take the Output."
	DescSharedModuleEnvVarNamespaceName  = "Name of the parent Namespace of the parent Module from which to take the Output."
	DescSharedModuleEnvVarOutputName     = "Name of Output to take as input."
	DescSharedModuleEnvVarSecretName     = "Name of the Secret to take as input."
	DescSharedModuleEnvVarSecretScope    = "Scope on which the Secret to take as input as been stored (Stack, Namespace or Module)"
)
