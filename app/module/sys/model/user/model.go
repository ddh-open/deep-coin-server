package user

import (
	"devops-http/app/module/base"
	uuid "github.com/satori/go.uuid"
)

// DevopsSysUser 用户
type DevopsSysUser struct {
	base.DevopsModel
	UUID      uuid.UUID `gorm:"column:uuid;type:varchar(64);default:null" json:"uuid"`                   // UUID
	Username  string    `gorm:"unique;column:username;type:varchar(32);unique;not null" json:"username"` // 登录名/11111
	Password  string    `gorm:"column:password;type:varchar(32);not null" json:"password"`               // 密码
	Salt      string    `gorm:"column:salt;type:varchar(16);not null;default:1111" json:"salt"`          // 密码盐
	Merchants string    `gorm:"column:merchants;type:varchar(64);" json:"merchants"`                     // 所属商户
	RealName  string    `gorm:"column:real_name;type:varchar(32);default:null" json:"realName"`          // 真实姓名
	UserType  int       `gorm:"column:user_type;type:int;default:null;default:1" json:"userType"`        // 用户类型
	Status    int       `gorm:"column:status;type:int;default:null;default:10" json:"status"`            // 状态
	WorkNum   string    `gorm:"column:work_num;type:varchar(20);default:null" json:"workNum"`            // 工号
	EndTime   string    `gorm:"column:end_time;type:varchar(32);default:null" json:"endTime"`            // 结束时间
	Email     string    `gorm:"column:email;type:varchar(64);default:null" json:"email"`                 // email
	Tel       string    `gorm:"column:tel;type:varchar(32);default:null" json:"tel"`                     // 手机号
	Address   string    `gorm:"column:address;type:varchar(32);default:null" json:"address"`             // 地址
	TitleURL  string    `gorm:"column:title_url;type:varchar(200);default:null" json:"titleURL"`         // 头像地址
	Remark    string    `gorm:"column:remark;type:varchar(1000);default:null" json:"remark"`             // 说明
	Theme     string    `gorm:"column:theme;type:varchar(64);default:null;default:default" json:"theme"` // 主题
	Enable    int       `gorm:"column:enable;type:tinyint(1);default:null;default:1" json:"enable"`      // 是否启用//radio/1,启用,2,禁用
	UpdateID  int64     `gorm:"column:update_id;type:bigint;default:null;default:0" json:"updateID"`     // 更新人
	CreateID  int64     `gorm:"column:create_id;type:bigint;default:null;default:0" json:"createID"`     // 创建者
}

// TableName 表名
func (s *DevopsSysUser) TableName() string {
	return "devops_sys_users"
}
