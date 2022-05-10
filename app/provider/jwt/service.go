package jwt

import (
	contract2 "devops-http/app/contract"
	"github.com/ddh-open/gin/framework/contract"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"time"
)

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

var ConcurrencyControl = &singleflight.Group{}

func NewJWTService(params ...interface{}) (interface{}, error) {
	config := params[1].(contract.Config)
	signingKey := config.GetString("app.jwt.signing-key")
	return &JWT{config: config, SigningKey: []byte(signingKey)}, nil
}

type JWT struct {
	SigningKey []byte
	config     contract.Config
}

func (j *JWT) CreateClaims(baseClaims contract2.BaseClaims) contract2.CustomClaims {
	claims := contract2.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(j.config.GetInt("app.jwt.buffer-time")), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                                           // 签名生效时间
			ExpiresAt: time.Now().Unix() + int64(j.config.GetInt("app.jwt.expires-time")), // 过期时间 7天  配置文件
			Issuer:    j.config.GetString("app.jwt.issuer"),                               // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims contract2.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims contract2.CustomClaims) (string, error) {
	v, err, _ := ConcurrencyControl.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*contract2.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &contract2.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*contract2.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
