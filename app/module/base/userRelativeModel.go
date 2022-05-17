package base

import uuid "github.com/satori/go.uuid"

const (
	MENUS = iota
	APIS
)

var SourceList = [...]string{"MENUS", "APIS"}

// DevopsSysUser 用户
type DevopsSysUser struct {
	DevopsModel
	UUID     uuid.UUID        `gorm:"column:uuid;type:varchar(64);default:null" json:"uuid"`                   // UUID
	Username string           `gorm:"unique;column:username;type:varchar(32);unique;not null" json:"username"` // 登录名/11111
	Password string           `gorm:"column:password;type:varchar(32);not null" json:"password"`               // 密码
	Salt     string           `gorm:"column:salt;type:varchar(16);not null;default:1111" json:"salt"`          // 密码盐
	Domain   string           `gorm:"column:domain;type:varchar(64);default:default" json:"merchants"`         // 所属域
	RealName string           `gorm:"column:real_name;type:varchar(32);default:null" json:"realName"`          // 真实姓名
	UserType int              `gorm:"column:user_type;type:int;default:null;default:1" json:"userType"`        // 用户类型
	Status   int              `gorm:"column:status;type:int;default:null;default:10" json:"status"`            // 状态
	WorkNum  string           `gorm:"column:work_num;type:varchar(20);default:null" json:"workNum"`            // 工号
	EndTime  string           `gorm:"column:end_time;type:varchar(32);default:null" json:"endTime"`            // 结束时间
	Email    string           `gorm:"column:email;type:varchar(64);default:null" json:"email"`                 // email
	Tel      string           `gorm:"column:tel;type:varchar(32);default:null" json:"tel"`                     // 手机号
	Address  string           `gorm:"column:address;type:varchar(32);default:null" json:"address"`             // 地址
	TitleURL string           `gorm:"column:title_url;type:varchar(200);default:null" json:"titleURL"`         // 头像地址
	Remark   string           `gorm:"column:remark;type:varchar(1000);default:null" json:"remark"`             // 说明
	Theme    string           `gorm:"column:theme;type:varchar(64);default:null;default:default" json:"theme"` // 主题
	Enable   bool             `gorm:"column:enable;default:true" json:"enable"`                                // 是否启用//radio/1,启用,2,禁用
	Groups   []DevopsSysGroup `json:"groups" gorm:"many2many:devops_sys_group_relative_users;"`
}

// TableName 表名
func (s *DevopsSysUser) TableName() string {
	return "devops_sys_users"
}

// DevopsSysGroup 分组表（组织架构）
type DevopsSysGroup struct {
	DevopsModel
	ParentID  int             `gorm:"column:parent_id;type:int;default:null;default:0" json:"parentId"`   // 上级机构
	Name      string          `gorm:"unique;column:name;type:varchar(32);not null" json:"name"`           // 部门/11111
	Code      string          `gorm:"column:code;type:varchar(128);default:null" json:"code"`             // 机构编码
	Sort      int             `gorm:"column:sort;type:int;default:null;default:0" json:"sort"`            // 序号
	Linkman   string          `gorm:"column:linkman;type:varchar(64);default:null" json:"linkman"`        // 联系人
	LinkmanNo string          `gorm:"column:linkman_no;type:varchar(32);default:null" json:"linkmanNo"`   // 联系人电话
	Remark    string          `gorm:"column:remark;type:varchar(128);default:null" json:"remark"`         // 组描述
	Enable    int             `gorm:"column:enable;type:tinyint(1);default:null;default:1" json:"enable"` // 是否启用
	Alias     string          `gorm:"column:alias;type:varchar(128);default:null" json:"alias"`           // 别名
	Wechat    string          `gorm:"column:wechat;type:varchar(128);default:null" json:"wechat"`         // wechat
	Domain    uint            `json:"domain" gorm:"column:domain;type:int;default:null"`                  // 域
	Users     []DevopsSysUser `json:"users" gorm:"many2many:devops_sys_group_relative_users;"`

	Children   []DevopsSysGroup `json:"children" gorm:"-"`
	ParentName string           `json:"parentName" gorm:"-"`
}
