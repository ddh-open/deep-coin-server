package contract

import (
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

// JWT 定义字符串凭证
const JWT = "user:jwt"

// CustomClaims Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	UUID     uuid.UUID
	ID       uint
	Username string
	NickName string
}

// JWTService 定义接口
type JWTService interface {
	// CreateClaims 创建
	CreateClaims(baseClaims BaseClaims) CustomClaims
	// CreateToken 创建一个token
	CreateToken(claims CustomClaims) (string, error)
	// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
	CreateTokenByOldToken(oldToken string, claims CustomClaims) (string, error)
	// ParseToken 解析 token
	ParseToken(tokenString string) (*CustomClaims, error)
}

//func Verify(st interface{}, roleMap Rules) (err error) {
//	compareMap := map[string]bool{
//		"lt": true,
//		"le": true,
//		"eq": true,
//		"ne": true,
//		"ge": true,
//		"gt": true,
//	}
//
//	typ := reflect.TypeOf(st)
//	val := reflect.ValueOf(st) // 获取reflect.Type类型
//
//	kd := val.Kind() // 获取到st对应的类别
//	if kd != reflect.Struct {
//		return errors.New("expect struct")
//	}
//	num := val.NumField()
//	// 遍历结构体的所有字段
//	for i := 0; i < num; i++ {
//		tagVal := typ.Field(i)
//		val := val.Field(i)
//		if len(roleMap[tagVal.Name]) > 0 {
//			for _, v := range roleMap[tagVal.Name] {
//				switch {
//				case v == "notEmpty":
//					if isBlank(val) {
//						return errors.New(tagVal.Name + "值不能为空")
//					}
//				case strings.Split(v, "=")[0] == "regexp":
//					if !regexpMatch(strings.Split(v, "=")[1], val.String()) {
//						return errors.New(tagVal.Name + "格式校验不通过")
//					}
//				case compareMap[strings.Split(v, "=")[0]]:
//					if !compareVerify(val, v) {
//						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
//					}
//				}
//			}
//		}
//	}
//	return nil
//}
