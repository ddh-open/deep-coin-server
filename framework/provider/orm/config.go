package orm

import (
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"gorm.io/gorm"
)

// GetBaseConfig 读取database.yaml根目录结构
func GetBaseConfig(c framework.Container) *contract2.DBConfig {
	configService := c.MustMake(contract2.ConfigKey).(contract2.Config)
	logService := c.MustMake(contract2.LogKey).(contract2.Log)
	config := &contract2.DBConfig{}
	// 直接使用配置服务的load方法读取,yaml文件
	err := configService.Load("database", config)
	if err != nil {
		// 直接使用logService来打印错误信息
		logService.Error("parse database config error")
		return nil
	}
	return config
}

// WithConfigPath 加载配置文件地址
func WithConfigPath(configPath string) contract2.DBOption {
	return func(container framework.Container, config *contract2.DBConfig) error {
		configService := container.MustMake(contract2.ConfigKey).(contract2.Config)
		// 加载configPath配置路径
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}

// WithGormConfig 表示自行配置Gorm的配置信息
func WithGormConfig(gormConfig *gorm.Config) contract2.DBOption {
	return func(container framework.Container, config *contract2.DBConfig) error {
		if gormConfig.Logger == nil {
			gormConfig.Logger = config.Logger
		}
		config.Config = gormConfig
		return nil
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() contract2.DBOption {
	return func(container framework.Container, config *contract2.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() contract2.DBOption {
	return func(container framework.Container, config *contract2.DBConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
