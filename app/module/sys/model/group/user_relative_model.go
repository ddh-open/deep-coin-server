package group

// DevopsSysGroupRelativeUser 用户和应用分组关系表
type DevopsSysGroupRelativeUser struct {
	DevopsSysGroupId uint `json:"group_id" gorm:"primaryKey"`
	DevopsSysUserId  uint `json:"user_id" gorm:"primaryKey"`
}

// TableName 表名
func (s *DevopsSysGroupRelativeUser) TableName() string {
	return "devops_sys_group_relative_users"
}
