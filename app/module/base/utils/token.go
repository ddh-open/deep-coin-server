package utils

import (
	"devops-http/app/module/base"
	"devops-http/framework/gin"
	"github.com/pkg/errors"
)

func ParseToken(c *gin.Context) (*base.TokenUser, error) {
	userInterface, flag := c.Get(c.GetUserKey())
	if !flag {
		return nil, errors.New("未获取token的user")
	}
	if userData, ok := userInterface.(*base.TokenUser); ok {
		return userData, nil
	} else {
		return nil, errors.New("用户信息解析失败...")
	}
}
