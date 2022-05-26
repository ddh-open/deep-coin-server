package application

import "devops-http/app/module/base"

// DevopsClusterApplication 服务结构体
type DevopsClusterApplication struct {
	base.DevopsModel
	Name    string                           `json:"name" gorm:"comment:服务名"`    // 服务名
	Remark  string                           `json:"remark" gorm:"comment:服务描述"` // 服务描述
	Groups  []base.DevopsSysGroup            `json:"groups" gorm:"many2many:devops_cluster_application_relative_groups;"`
	Configs []DevopsClusterApplicationConfig `json:"configs" gorm:"many2many:devops_cluster_application_relative_configs;"`
}

type DevopsClusterApplicationRelativeGroup struct {
	DevopsClusterApplicationId uint `json:"devopsClusterApplicationId" gorm:"primaryKey"`
	DevopsSysGroupId           uint `json:"devopsSysGroupId" gorm:"primaryKey"`
}

type DevopsClusterApplicationRelativeConfig struct {
	DevopsClusterApplicationId       uint `json:"devopsClusterApplicationId" gorm:"primaryKey"`
	DevopsClusterApplicationConfigId uint `json:"DevopsClusterApplicationConfigId" gorm:"primaryKey"`
}

// DevopsClusterApplicationConfig 服务配置结构体
type DevopsClusterApplicationConfig struct {
	base.DevopsModel
	Name      string `json:"name" gorm:"comment:服务名"`       // 配置名
	Env       string `json:"env" gorm:"env"`                // 环境名
	Remark    string `json:"remark" gorm:"comment:服务环境描述"`  // 服务环境描述
	Namespace string `json:"namespace" gorm:"comment:命令空间"` // 命令空间
}
