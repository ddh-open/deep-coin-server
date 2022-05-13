package config

import (
	framework2 "devops-http/framework"
	contract2 "devops-http/framework/contract"
	"path/filepath"
)

type NiceConfigProvider struct{}

// Register register a new function for make a service instance
func (provider *NiceConfigProvider) Register(c framework2.Container) framework2.NewInstance {
	return NewNiceConfig
}

// Boot will called when the service instantiate
func (provider *NiceConfigProvider) Boot(c framework2.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *NiceConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *NiceConfigProvider) Params(c framework2.Container) []interface{} {
	appService := c.MustMake(contract2.AppKey).(contract2.App)
	envService := c.MustMake(contract2.EnvKey).(contract2.Env)
	env := envService.AppEnv()
	// 配置文件夹地址
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Name / Name define the name for this service
func (provider *NiceConfigProvider) Name() string {
	return contract2.ConfigKey
}
