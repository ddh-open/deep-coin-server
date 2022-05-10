package ldap

import (
	"devops-http/app/contract"
	"github.com/ddh-open/gin/framework"
	contract2 "github.com/ddh-open/gin/framework/contract"
)

// ProviderLdap 提供ldap的具体实现方法
type ProviderLdap struct {
	Config contract2.Config
}

// Register 注册App方法
func (h *ProviderLdap) Register(container framework.Container) framework.NewInstance {
	return NewLdapService
}

// Boot 启动调用
func (h *ProviderLdap) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *ProviderLdap) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *ProviderLdap) Params(container framework.Container) []interface{} {
	// 获取配置参数
	h.Config = container.MustMake(contract2.ConfigKey).(contract2.Config)
	return []interface{}{container, h.Config}
}

// Name 获取字符串凭证
func (h *ProviderLdap) Name() string {
	return contract.KeyLdap
}
