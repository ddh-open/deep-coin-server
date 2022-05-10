package contract

// KeyLdap 定义字符串凭证
const KeyLdap = "user:ldap"

// Ldap 定义接口
type Ldap interface {
	Search(username string, filter string, keys ...string) (map[string]string, error)
	Login(username string, password string) (map[string]string, error)
	ChangePassword(username string, password string) error
	CreateOu(ou string, baseDN string, attr map[string][]string) (err error)
	DeleteOu(ou string, baseDN string) (err error)
	CreateUser(cn string, baseDN string, password string, attr map[string][]string) (err error)
	DeleteUser(cn string, baseDN string) (err error)
	Init() error
	Close()
}
