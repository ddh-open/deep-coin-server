package cache

import (
	framework2 "devops-http/framework"
	contract2 "devops-http/framework/contract"
	services2 "devops-http/framework/provider/cache/services"
	"strings"
)

// NiceCacheProvider 服务提供者
type NiceCacheProvider struct {
	framework2.ServiceProvider

	Driver string // Driver
}

// Register 注册一个服务实例
func (l *NiceCacheProvider) Register(c framework2.Container) framework2.NewInstance {
	if l.Driver == "" {
		tcs, err := c.Make(contract2.ConfigKey)
		if err != nil {
			// 默认使用console
			return services2.NewMemoryCache
		}

		cs := tcs.(contract2.Config)
		l.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	// 根据driver的配置项确定
	switch l.Driver {
	case "redis":
		return services2.NewRedisCache
	case "memory":
		return services2.NewMemoryCache
	default:
		return services2.NewMemoryCache
	}
}

// Boot 启动的时候注入
func (l *NiceCacheProvider) Boot(c framework2.Container) error {
	return nil
}

// IsDefer 是否延迟加载
func (l *NiceCacheProvider) IsDefer() bool {
	return true
}

// Params 定义要传递给实例化方法的参数
func (l *NiceCacheProvider) Params(c framework2.Container) []interface{} {
	return []interface{}{c}
}

// Name 定义对应的服务字符串凭证
func (l *NiceCacheProvider) Name() string {
	return contract2.CacheKey
}
