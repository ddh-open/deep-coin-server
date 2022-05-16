package middleware

import (
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
	"strings"
)

var whiteList = append(make([]string, 0), "/user/login", "/user/register")

// Auth 鉴权中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		authPass := false
		for i := range whiteList {
			if c.Request.RequestURI == whiteList[i] {
				authPass = true
			}
		}
		// debug 模式开启， 跳过权限验证
		if c.MustMakeConfig().GetBool("app.debug") {
			authPass = true
		}
		// 不在白名单， 通过grpc鉴权
		if !authPass {
			if !strings.HasPrefix(token, "Bearer ") {
				response.FailWithDetailed(gin.H{"reload": true}, "token 格式不正确", c)
				c.Abort()
				return
			}
			token = strings.Replace(token, "Bearer ", "", 1)
			jwt := c.MustMake(contract.JWT).(contract.JWTService)
			tokenInfo, err := jwt.ParseToken(token)
			if err != nil {
				response.FailWithDetailed(gin.H{"reload": true}, "解析token出错:", c)
				c.Abort()
				return
			}
			domain := c.Request.Header.Get("domain")
			userInfo := parseUser(tokenInfo, domain)
			// 鉴权

			c.Set(c.GetUserKey(), userInfo)
		}
		c.Next()
	}
}

func parseUser(data *contract.CustomClaims, domain string) *base.TokenUser {
	if domain == "" {
		domain = "default"
	}
	return &base.TokenUser{
		Id:            int(data.ID),
		Uuid:          data.UUID.String(),
		Username:      data.Username,
		RealName:      data.NickName,
		CurrentDomain: domain,
	}
}
