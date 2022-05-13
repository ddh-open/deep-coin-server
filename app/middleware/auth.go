package middleware

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/response"
	"devops-http/framework/gin"
	"devops-http/resources/proto/userGrpc"
)

var whiteList = append(make([]string, 0), "/user/login", "/user/register")

// Auth 登录中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("bearer")
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
			devops := c.MustMake(contract.KeyGrpc).(contract.ServiceGrpc)
			conn, err := devops.GetGrpc("grpc.user")
			if err != nil {
				response.FailWithDetailed(gin.H{"reload": true}, "用户api鉴权连接出错:"+err.Error(), c)
				c.Abort()
				return
			}
			defer conn.Close()
			client := userGrpc.NewServiceAuthClient(conn)
			resp, err := client.AuthApi(context.Background(), &userGrpc.AuthApiRequest{
				Path:  c.Request.RequestURI,
				Token: token,
			})
			if err != nil {
				response.FailWithDetailed(gin.H{"reload": true}, "用户api鉴权调用出错:"+err.Error(), c)
				c.Abort()
				return
			}
			// 代表响应成功
			if resp.GetResult().GetCode() != 200 {
				response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问:"+resp.GetResult().GetMsg(), c)
				c.Abort()
				return
			}
			domain := c.Request.Header.Get("domain")
			userInfo := parseUser(resp.GetUser(), domain)
			c.Set(c.GetUserKey(), userInfo)
		}
		c.Next()
	}
}

func parseUser(data interface{}, domain string) *base.TokenUser {
	if domain == "" {
		domain = "freemud"
	}
	if userData, ok := data.(*userGrpc.BaseUserInfo); ok {
		return &base.TokenUser{
			Id:            int(userData.Id),
			Uuid:          userData.Uuid,
			Username:      userData.Username,
			RealName:      userData.Nickname,
			CurrentDomain: domain,
		}
	}
	return nil
}
