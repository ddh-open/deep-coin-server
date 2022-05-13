package app

import (
	framework2 "devops-http/framework"
	"devops-http/framework/contract"
)

// NiceAppProvider 提供App的具体实现方法
type NiceAppProvider struct {
	BaseFolder string
}

// Register 注册NiceApp方法
func (h *NiceAppProvider) Register(container framework2.Container) framework2.NewInstance {
	return NewNiceApp
}

// Boot 启动调用
func (h *NiceAppProvider) Boot(container framework2.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *NiceAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (h *NiceAppProvider) Params(container framework2.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name 获取字符串凭证
func (h *NiceAppProvider) Name() string {
	return contract.AppKey
}
