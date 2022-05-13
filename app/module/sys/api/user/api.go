package user

import (
	"devops-http/app/module/sys/service/user"
	"devops-http/framework"
	"devops-http/framework/gin"
)

type ApiUser struct {
	service *user.Service
}

func Register(r *gin.Engine) error {
	api := NewUserApi(r.GetContainer())
	userGroup := r.Group("/user/", func(c *gin.Context) {
	})
	userGroup.POST("login", api.Login)                   // 用户登录
	userGroup.POST("changePassword", api.ChangePassword) // 更改用户名密码
	userGroup.PUT("modify", api.Modify)                  // 修改用户的信息
	userGroup.POST("logout", api.Logout)                 // 用户退出
	userGroup.GET("info", api.GetUserInfo)               // 用户详情
	userGroup.POST("add", api.Add)                       // 用户新增
	userGroup.DELETE("delete", api.Delete)               // 用户删除
	userGroup.POST("list", api.UserList)                 // 用户列表
	userGroup.POST("register", api.Register)             // 用户注册

	// 用户关联角色相关
	userGroup.GET("relative/roles/:id", api.UserRelativeRole)             // 获取用户的相关权限
	userGroup.GET("relative/apis", api.UserGetApis)                       // 获取用户的相关权限
	userGroup.POST("relative/roles/add", api.UserRelativeRoleAdd)         // 给用户增加角色
	userGroup.DELETE("relative/roles/delete", api.UserRelativeRoleDelete) // 给用户删除角色
	return nil
}

func NewUserApi(c framework.Container) *ApiUser {
	return &ApiUser{service: user.NewService(c)}
}
