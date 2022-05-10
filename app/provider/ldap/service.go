package ldap

import (
	"crypto/tls"
	"fmt"
	"github.com/ddh-open/gin/framework/contract"
	"github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"time"
)

func NewLdapService(params ...interface{}) (interface{}, error) {
	return &Ldap{config: params[1].(contract.Config), ldap: nil}, nil
}

type Ldap struct {
	ldap   *ldap.Conn
	config contract.Config
}

func (l *Ldap) Init() error {
	var conn *ldap.Conn
	var err error
	if l.config.GetString("app.ldap.port") == "636" {
		conn, err = ldap.DialTLS("tcp", l.config.GetString("app.ldap.host")+":"+l.config.GetString("app.ldap.port"), &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", l.config.GetString("app.ldap.host")+":"+l.config.GetString("app.ldap.port"))
	}
	l.ldap = conn
	if l.ldap != nil {
		l.ldap.SetTimeout(5 * time.Second)
	}
	return err
}

func (l *Ldap) Close() {
	if l.ldap != nil {
		l.ldap.Close()
		l.ldap = nil
	}
}

func (l *Ldap) Search(username string, filter string, keys ...string) (map[string]string, error) {
	if filter == "" {
		filter = fmt.Sprintf(l.config.GetString("app.ldap.search_filter"), username)
	}
	var cur *ldap.SearchResult
	err := l.ldap.Bind(l.config.GetString("app.ldap.bind_dn"), l.config.GetString("app.ldap.bind_passwd"))
	if err != nil {
		return nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		l.config.GetString("app.ldap.search_ou"),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", "cn", "mail", "displayName"},
		nil,
	)

	if cur, err = l.ldap.Search(searchRequest); err != nil {
		return nil, err
	}

	if len(cur.Entries) == 0 {
		return nil, nil
	}
	result := map[string]string{
		"dn":   cur.Entries[0].DN,
		"cn":   cur.Entries[0].GetAttributeValue("cn"),
		"mail": cur.Entries[0].GetAttributeValue("mail"),
	}
	return result, nil
}

func (l *Ldap) Login(username string, password string) (map[string]string, error) {
	err := l.Init()
	if err != nil {
		return nil, err
	}
	defer l.Close()
	entry, err := l.Search(username, "")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("该用户不存在")
	}
	err = l.ldap.Bind(entry["dn"], password)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (l *Ldap) ChangePassword(username string, password string) (err error) {
	err = l.Init()
	if err != nil {
		return
	}
	defer l.Close()
	// 登录管理员账户
	entry, err := l.Search(username, "")
	if err != nil {
		return
	}
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("%q", password))
	if err != nil {
		return
	}
	modReq := ldap.NewModifyRequest(entry["dn"], []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})
	err = l.ldap.Modify(modReq)
	return
}

func (l *Ldap) CreateUser(cn string, baseDN string, password string, attr map[string][]string) (err error) {
	if password == "" {
		password = "1234!@#$.com"
	}
	err = l.Init()
	if err != nil {
		return
	}
	defer l.Close()
	// 登录管理员账户
	err = l.ldap.Bind(l.config.GetString("app.ldap.bind_dn"), l.config.GetString("app.ldap.bind_passwd"))
	if err != nil {
		return err
	}

	// 创建用户
	request := ldap.NewAddRequest(fmt.Sprintf("CN=%s,%s", cn, baseDN), []ldap.Control{})
	request.Attribute("name", []string{cn})
	request.Attribute("objectClass", []string{"top", "organizationalPerson", "user", "person"})
	request.Attribute("sAMAccountName", []string{cn})
	request.Attribute("userPrincipalName", []string{fmt.Sprintf("%v@office.freemud.cn", cn)})
	if attr != nil {
		for key, val := range attr {
			request.Attribute(key, val)
		}
	}
	err = l.ldap.Add(request)
	if err != nil {
		return
	}

	// 修改密码
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, err := utf16.NewEncoder().String(fmt.Sprintf("%q", password))
	if err != nil {
		return
	}
	modReq := ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", cn, baseDN), []ldap.Control{})
	modReq.Replace("unicodePwd", []string{pwdEncoded})

	err = l.ldap.Modify(modReq)
	if err != nil {
		return
	}

	modReq = ldap.NewModifyRequest(fmt.Sprintf("CN=%s,%s", cn, baseDN), []ldap.Control{})
	// 解除账户禁用
	modReq.Replace("userAccountControl", []string{"544"})
	err = l.ldap.Modify(modReq)
	return
}

func (l *Ldap) DeleteUser(cn string, baseDN string) (err error) {
	err = l.Init()
	if err != nil {
		return
	}
	defer l.Close()
	// 登录管理员账户
	err = l.ldap.Bind(l.config.GetString("app.ldap.bind_dn"), l.config.GetString("app.ldap.bind_passwd"))
	if err != nil {
		return err
	}

	request := ldap.NewDelRequest(fmt.Sprintf("CN=%s,%s", cn, baseDN), nil)
	err = l.ldap.Del(request)
	return
}

func (l *Ldap) CreateOu(ou string, baseDN string, attr map[string][]string) (err error) {
	err = l.Init()
	if err != nil {
		return
	}
	defer l.Close()
	// 登录管理员账户
	err = l.ldap.Bind(l.config.GetString("app.ldap.bind_dn"), l.config.GetString("app.ldap.bind_passwd"))
	if err != nil {
		return err
	}
	request := ldap.NewAddRequest(fmt.Sprintf("OU=%s,%s", ou, baseDN), []ldap.Control{})
	request.Attribute("ou", []string{ou})
	request.Attribute("objectClass", []string{"top", "organizationalUnit"})
	if attr != nil {
		for key, val := range attr {
			request.Attribute(key, val)
		}
	}
	err = l.ldap.Add(request)
	return
}

func (l *Ldap) DeleteOu(ou string, baseDN string) (err error) {
	err = l.Init()
	if err != nil {
		return
	}
	defer l.Close()
	// 登录管理员账户
	err = l.ldap.Bind(l.config.GetString("app.ldap.bind_dn"), l.config.GetString("app.ldap.bind_passwd"))
	if err != nil {
		return err
	}
	request := ldap.NewDelRequest(fmt.Sprintf("OU=%s,%s", ou, baseDN), nil)
	err = l.ldap.Del(request)
	return
}
