package host

import "devops-http/app/module/base"

// DevopsCmdbHost cmdb 主机列表
type DevopsCmdbHost struct {
	base.DevopsModel
	Ip            string                `json:"ip" gorm:"comment:主机Ip"`           // 主机Ip
	Port          string                `json:"port" gorm:"comment:主机ssh端口"`      // 主机ssh端口
	Name          string                `json:"name" gorm:"comment:主机名"`          // 主机名
	ConfigureInfo string                `json:"configureInfo" gorm:"comment:主机名"` // 配置信息
	SystemInfo    string                `json:"systemInfo" gorm:"comment:系统信息"`   // 系统信息
	Status        string                `json:"status" gorm:"comment:主机状态"`       // 主机的状态
	Remark        string                `json:"remark" gorm:"comment:主机描述"`       // 主机描述
	Checked       bool                  `json:"checked" gorm:"comment:是否验证过"`     // 是否验证过
	Sort          int                   `json:"sort" gorm:"comment:排序标记"`         // 排序标记
	Groups        []DevopsCmdbHostGroup `json:"groups" gorm:"many2many:devops_cmdb_host_group_relative_hosts;"`
}

type DevopsCmdbHostGroup struct {
	base.DevopsModel
	ParentId uint                  `json:"parentId" gorm:"default:0"`
	Name     string                `json:"name" gorm:"comment:分组名"`                 // 分组名
	Remark   string                `json:"remark" gorm:"comment:分组的描述"`             // 分组的描述
	Enable   bool                  `json:"enable" gorm:"comment:是否启用;default:true"` // 分组状态
	HostNum  int                   `json:"hostNum" gorm:"comment:主机数量;default:0"`   // 主机数量
	Sort     int                   `json:"sort" gorm:"comment:排序;default:0"`        // 排序
	Children []DevopsCmdbHostGroup `json:"children" gorm:"-"`                       // 主机数量
	Hosts    []DevopsCmdbHost      `json:"hosts"  gorm:"many2many:devops_cmdb_host_group_relative_hosts;"`
}

// DevopsCmdbHostGroupRelativeHost 主机和分组关系表
type DevopsCmdbHostGroupRelativeHost struct {
	DevopsCmdbHostGroupId uint `json:"devopsCmdbHostGroupId" gorm:"primaryKey"`
	DevopsCmdbHostId      uint `json:"devops_cmdb_host_id" gorm:"primaryKey"`
}
