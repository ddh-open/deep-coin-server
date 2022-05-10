package domain

import "devops-http/app/module/base"

// DevopsSysDomain 域表
type DevopsSysDomain struct {
	base.DevopsModel
	Name        string `gorm:"column:name;type:varchar(256);unique;not null" json:"name"`            // 域账户名
	DomainNum   string `gorm:"column:domain_num;type:varchar(256);unique;not null" json:"domainNum"` // 商户号
	EnglishName string `gorm:"column:english_name;type:varchar(256);unique;" json:"englishName"`     // 英文名
	Remark      string `gorm:"column:remark;type:varchar(256);" json:"remark"`                       // 描述
	Enable      int    `gorm:"column:enable;type:tinyint(1);default:null;default:1" json:"enable"`   // 是否启用
	UpdateID    int64  `gorm:"column:update_id;type:bigint;default:null;default:0" json:"updateId"`  // 更新人
	CreateID    int64  `gorm:"column:create_id;type:bigint;default:null;default:0" json:"createId"`  // 创建者
}
