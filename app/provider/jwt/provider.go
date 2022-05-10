package jwt

import (
	"devops-http/app/contract"
	"github.com/ddh-open/gin/framework"
	contract2 "github.com/ddh-open/gin/framework/contract"
)

// ProviderJWT 提供jwt的具体实现方法
type ProviderJWT struct {
	Config contract2.Config
}

// Register 注册App方法
func (h *ProviderJWT) Register(container framework.Container) framework.NewInstance {
	return NewJWTService
}

// Boot 启动调用
func (h *ProviderJWT) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *ProviderJWT) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *ProviderJWT) Params(container framework.Container) []interface{} {
	// 获取配置参数
	h.Config = container.MustMake(contract2.ConfigKey).(contract2.Config)
	return []interface{}{container, h.Config}
}

// Name 获取字符串凭证
func (h *ProviderJWT) Name() string {
	return contract.JWT
}
