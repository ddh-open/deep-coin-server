package gin

import (
	contract2 "devops-http/framework/contract"
)

var userKey = "claims@user"

func (c *Context) GetUserKey() string {
	return userKey
}

// MustMakeApp 从容器中获取App服务
func (c *Context) MustMakeApp() contract2.App {
	return c.MustMake(contract2.AppKey).(contract2.App)
}

// MustMakeKernel 从容器中获取Kernel服务
func (c *Context) MustMakeKernel() contract2.Kernel {
	return c.MustMake(contract2.KernelKey).(contract2.Kernel)
}

// MustMakeConfig 从容器中获取配置服务
func (c *Context) MustMakeConfig() contract2.Config {
	return c.MustMake(contract2.ConfigKey).(contract2.Config)
}

// MustMakeLog 从容器中获取日志服务
func (c *Context) MustMakeLog() contract2.Log {
	return c.MustMake(contract2.LogKey).(contract2.Log)
}
