package orm

import (
	framework2 "devops-http/framework"
	"devops-http/framework/contract"
)

// GormProvider 提供App的具体实现方法
type GormProvider struct {
	defaultPath string
}

// Register 注册方法
func (h *GormProvider) Register(container framework2.Container) framework2.NewInstance {
	return NewNiceGorm
}

// Boot 启动调用
func (h *GormProvider) Boot(container framework2.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *GormProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (h *GormProvider) Params(container framework2.Container) []interface{} {
	defaultPath := "database.mysql"
	if h.defaultPath != "" {
		defaultPath = h.defaultPath
	}
	return []interface{}{container, defaultPath}
}

// Name 获取字符串凭证
func (h *GormProvider) Name() string {
	return contract.ORMKey
}
