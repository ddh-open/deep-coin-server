package env

import (
	framework2 "devops-http/framework"
	contract2 "devops-http/framework/contract"
)

type NiceEnvProvider struct {
	Folder string
}

// Register register a new function for make a service instance
func (provider *NiceEnvProvider) Register(c framework2.Container) framework2.NewInstance {
	return NewNiceEnv
}

// Boot will called when the service instantiate
func (provider *NiceEnvProvider) Boot(c framework2.Container) error {
	app := c.MustMake(contract2.AppKey).(contract2.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NiceEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NiceEnvProvider) Params(c framework2.Container) []interface{} {
	return []interface{}{provider.Folder}
}

// Name / Name define the name for this service
func (provider *NiceEnvProvider) Name() string {
	return contract2.EnvKey
}
