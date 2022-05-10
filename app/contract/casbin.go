package contract

import (
	"github.com/casbin/casbin/v2"
)

// KeyCaBin 定义字符串凭证
const KeyCaBin = "user:cabin"

// Cabin 定义接口
type Cabin interface {
	GetCabin() *casbin.SyncedEnforcer
	ClearCabin(v int, p ...string) bool
	UpdateCabinApi(oldPath string, newPath string, oldMethod string, newMethod string) error
}
