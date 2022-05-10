package base

import (
	"devops-http/app/module/base/dtime"
	"gorm.io/gorm"
)

type TokenUser struct {
	Username      string   `json:"username"`
	RealName      string   `json:"realName"`
	Id            int      `json:"id"`
	Uuid          string   `json:"uuid"`
	BelongDomain  string   `json:"belongDomain"`
	HaveDomains   []string `json:"haveDomains"`
	CurrentDomain string   `json:"currentDomain"`
}

type DevopsColumns struct {
	Label string `json:"label"`
	Width string `json:"width"`
	Prop  string `json:"prop"`
}

type DevopsModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;autoIncrementIncrement;column:id;type:bigint;not null" json:"id"` // 主键
	CreatedAt dtime.Time     `json:"created_at"`
	UpdatedAt dtime.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
