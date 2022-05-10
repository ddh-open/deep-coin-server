package config

import "devops-http/app/module/base"

type DevopsSysConfig struct {
	base.DevopsModel
	Name     string `gorm:"column:name;type:varchar(256);unique;not null" json:"name"`           // 配置名
	Key      string `gorm:"column:key;type:varchar(256);unique;not null" json:"key"`             //  配置的key
	Value    string `gorm:"column:value;type:varchar(256);unique;not null" json:"value"`         //  配置的value
	Remark   string `gorm:"column:remark;type:varchar(256);" json:"remark"`                      // 描述
	UpdateID int64  `gorm:"column:update_id;type:bigint;default:null;default:0" json:"updateId"` // 更新人
	CreateID int64  `gorm:"column:create_id;type:bigint;default:null;default:0" json:"createId"` // 创建者
}
