package casbin

import (
	"devops-http/app/contract"
	"github.com/ddh-open/gin/framework"
)

// ProviderCabin 提供tencent的具体实现方法
type ProviderCabin struct {
	ResourcePath string
}

// Register 注册NiceApp方法
func (h *ProviderCabin) Register(container framework.Container) framework.NewInstance {
	return NewCaBinService
}

// Boot 启动调用
func (h *ProviderCabin) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *ProviderCabin) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *ProviderCabin) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.ResourcePath}
}

// Name 获取字符串凭证
func (h *ProviderCabin) Name() string {
	return contract.KeyCaBin
}
