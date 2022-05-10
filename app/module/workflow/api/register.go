package api

import (
	"devops-http/app/module/workflow/api/identity"
	"devops-http/app/module/workflow/api/proc"
	"devops-http/app/module/workflow/api/process"
	"devops-http/app/module/workflow/api/task"
	"github.com/ddh-open/gin/framework/gin"
)

func Register(r *gin.Engine) error {
	procApi := proc.NewDefProc(r.GetContainer())
	userGroup := r.Group("v1/workflow/proc", func(c *gin.Context) {
	})
	userGroup.POST("save", procApi.Save)           // 定义工作流接口
	userGroup.DELETE("delete/:id", procApi.Delete) // 删除定义的工作流接口
	userGroup.POST("list", procApi.List)           // 查询工作流接口

	processApi := process.NewInstProcess(r.GetContainer())
	processGroup := r.Group("v1/workflow/process", func(c *gin.Context) {
	})
	processGroup.POST("start", processApi.Start)                    // 开启一个工作流
	processGroup.POST("inst/notify", processApi.FindProcNotifyInst) // 获取抄送的实例
	processGroup.POST("inst/myself", processApi.FindProcInstMyself) // 我发起的实例
	processGroup.POST("inst", processApi.FindProcInstMyself)        // 待我处理的实例

	taskApi := task.NewTaskApi(r.GetContainer())
	taskGroup := r.Group("v1/workflow/task", func(c *gin.Context) {
	})
	taskGroup.POST("complete", taskApi.CompleteTask) // 审批
	taskGroup.POST("draw", taskApi.WithDrawTask)     // 撤销

	identityApi := identity.NewApiIdentity(r.GetContainer())
	identityGroup := r.Group("v1/workflow/identity", func(c *gin.Context) {
	})
	identityGroup.POST(":id", identityApi.GetParticipant)
	return nil
}
