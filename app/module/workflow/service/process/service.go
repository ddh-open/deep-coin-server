package process

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/base/workflow"
	"devops-http/app/module/base/workflow/flow"
	"devops-http/app/module/workflow/model/execution"
	"devops-http/app/module/workflow/model/identity"
	"devops-http/app/module/workflow/model/node"
	"devops-http/app/module/workflow/model/proc"
	"devops-http/app/module/workflow/model/process"
	"devops-http/app/module/workflow/model/task"
	nodeSvc "devops-http/app/module/workflow/service/node"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Service struct {
	repository *base.Repository
}

func NewService(c framework.Container) *Service {
	db, err := c.MustMake(contract2.ORMKey).(contract2.ORMService).GetDB()
	logger := c.MustMake(contract2.LogKey).(contract2.Log)
	db.AutoMigrate(&process.WorkflowInstProc{})
	if err != nil {
		logger.Error("service 获取db出错： err", zap.Error(err))
	}
	return &Service{base.NewRepository(db)}
}

// StartProcess 流程启动
func (s *Service) StartProcess(req *request.ReceiverProcess) (id uint, err error) {
	//request.Company = userinfo.Company
	//request.Department = userinfo.Department
	//request.UserID = userinfo.ID
	//request.Username = userinfo.Username

	//if len(userinfo.Company) == 0 {
	//	return 0, errors.New("保存在redis中的【用户信息 userinfo】字段 company 不能为空")
	//}
	//if len(userinfo.Username) == 0 {
	//	return 0, errors.New("保存在redis中的【用户信息 userinfo】字段 username 不能为空")
	//}
	//if len(userinfo.ID) == 0 {
	//	return 0, errors.New("保存在redis中的【用户信息 userinfo】字段 ID 不能为空")
	//}
	//if len(userinfo.Department) == 0 {
	//	return 0, errors.New("保存在redis中的【用户信息 userinfo】字段 department 不能为空")
	//}

	// 获取流程定义
	procDef := make([]proc.WorkflowProc, 0)
	err = s.repository.GetDB().Model(&proc.WorkflowProc{}).Where(" name = ? and company = ? ", req.ProcName, req.Company).Find(&procDef).Error
	if err != nil {
		return 0, err
	}
	if len(procDef) <= 0 {
		return 0, errors.New("未找到定位的流程")
	}
	procData := procDef[0]
	nodeData := &node.Node{}
	err = json.Unmarshal([]byte(procData.Resource), nodeData)
	if err != nil {
		return 0, err
	}
	//--------以下需要添加事务-----------------
	step := 0 // 0 为开始节点
	// 新建流程实例
	var procInst = process.WorkflowInstProc{
		ProcDefID:     procData.ID,
		ProcDefName:   procData.Name,
		Title:         req.Title,
		Department:    req.Department,
		StartUserID:   req.UserID,
		StartUserName: req.Username,
		Company:       req.Company,
	}
	//开启事务
	tx := s.repository.GetDB().Begin()
	// 保存 实例
	err = tx.Create(&procInst).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id = procInst.ID
	exec := &execution.WorkflowExecution{
		ProcDefID:  procData.ID,
		ProcInstID: procInst.ID,
	}
	taskInst := &task.WorkflowTask{
		NodeID:        "start",
		ProcInstID:    procInst.ID,
		Assignee:      req.UserID,
		IsFinished:    true,
		ClaimTime:     utils.FormatDate(time.Now()),
		Step:          step,
		MemberCount:   1,
		UnCompleteNum: 0,
		ActType:       "or",
		AgreeNum:      1,
	}
	// 生成执行流，一串运行节点
	_, err = s.generateExec(exec, nodeData, req.UserID, req.Var, tx) //事务
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// 获取执行流信息
	var nodeInfos []*node.WorkflowNodeInfo
	err = json.Unmarshal([]byte(exec.NodeInfos), &nodeInfos)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// -----------------生成新任务-------------
	if nodeInfos[0].ActType == "and" {
		taskInst.UnCompleteNum = nodeInfos[0].MemberCount
		taskInst.MemberCount = nodeInfos[0].MemberCount
	}
	_, err = s.saveTaskByTx(taskInst, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//--------------------流转------------------
	// 流程移动到下一环节
	// 添加上一步的参与人
	workflowIdentityData := &identity.WorkflowIdentityLink{
		Type:       identity.KindsIdentity[identity.PARTICIPANT],
		UserID:     req.UserID,
		UserName:   req.Username,
		ProcInstID: id,
		Step:       step,
		Company:    req.Company,
		TaskID:     taskInst.ID,
		Comment:    "启动流程",
	}
	err = flow.MoveStage(nodeInfos, workflowIdentityData, "", true, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// fmt.Printf("流转到下一流程耗时：%v", time.Since(times))
	// fmt.Println("--------------提交事务----------")
	tx.Commit() //结束事务
	return
}

// GenerateExec 根据流程定义node生成执行流
func (s *Service) generateExec(e *execution.WorkflowExecution, n *node.Node, userID string, variable *map[string]string, tx *gorm.DB) (uint, error) {
	list, err := nodeSvc.ParseProcessConfig(n, variable)
	if err != nil {
		return 0, err
	}
	list.PushBack(node.WorkflowNodeInfo{
		NodeID: "结束",
		Type:   node.WorkflowNodeTypesList[node.START].Type,
	})
	list.PushFront(node.WorkflowNodeInfo{
		NodeID:   "开始",
		Type:     node.WorkflowNodeTypesList[node.START].Type,
		Approved: userID,
	})
	arr := utils.ListToArray(list)
	str, err := json.Marshal(&arr)
	if err != nil {
		return 0, err
	}
	e.NodeInfos = string(str)
	err = tx.Model(&execution.WorkflowExecution{}).Save(e).Error
	if err != nil {
		return 0, err
	}
	return e.ID, err
}

// saveTaskByTx 根据事物保存task任务到数据库
func (s *Service) saveTaskByTx(taskInst *task.WorkflowTask, tx *gorm.DB) (uint, error) {
	err := tx.Model(&task.WorkflowTask{}).Save(taskInst).Error
	if err != nil {
		return 0, err
	}
	return taskInst.ID, err
}

func (s *Service) FindProcList(req *workflow.ProcessPageReceiver) (result response.PageResult, err error) {
	var data []*process.WorkflowInstProc
	db := s.repository.GetDB().Model(&process.WorkflowInstProc{})
	// 流程的标题
	if req.Title != "" {
		db.Where("title like ?", "%"+req.Title+"%")
	}
	// 流程定义的名字
	if req.ProcDefName != "" {
		db.Where("proc_def_name like ?", "%"+req.ProcDefName+"%")
	}
	// 公司
	if req.Company != "" {
		db.Where("company like ?", "%"+req.Company+"%")
	}
	if req.ID != 0 {
		db.Where("id = ?", req.ID)
	}
	// 审核人
	if req.Candidate != "" {
		db.Where("candidate = ?", req.Candidate)
	}
	// 查询流程实例的列表
	err = db.Where("is_finished = 0").Count(&result.Total).Error
	if err != nil {
		return
	}
	err = db.Where("is_finished = 0").Offset(int(req.PageSize * req.Page)).Limit(int(req.PageSize)).Order("start_time desc").Find(&data).Error
	if err != nil {
		return
	}
	result.List = data
	result.PageSize = req.PageSize
	result.Page = req.Page
	return
}

// FindProcNotify 查询抄送我的流程
func (s *Service) FindProcNotify(req *workflow.ProcessPageReceiver) (result response.PageResult, err error) {
	var data []*process.WorkflowInstProc
	var ids []string
	db := s.repository.GetDB().Model(&identity.WorkflowIdentityLink{})
	err = db.Select("pro_inst_id").Where("type = notifier and company= ? and user_id = ?", req.Company, req.StartUserID).Find(&ids).Error
	if err != nil {
		return
	}
	db = s.repository.GetDB().Model(&process.WorkflowInstProc{})
	// 流程的标题
	if req.Title != "" {
		db.Where("title like ?", "%"+req.Title+"%")
	}
	// 流程定义的名字
	if req.ProcDefName != "" {
		db.Where("proc_def_name like ?", "%"+req.ProcDefName+"%")
	}
	// 公司
	if req.Company != "" {
		db.Where("company like ?", "%"+req.Company+"%")
	}
	if req.ID != 0 {
		db.Where("id = ?", req.ID)
	}
	db.Where("id in ?", strings.Join(ids, ","))
	// 查询流程实例的列表
	err = db.Count(&result.Total).Error
	if err != nil {
		return
	}
	err = db.Offset(int(req.PageSize * req.Page)).Limit(int(req.PageSize)).Order("start_time desc").Find(&data).Error
	return
}

// FindProcNotifyList 查找抄送我的流程实例
func (s *Service) FindProcNotifyList(req *workflow.ProcessPageReceiver) (result response.PageResult, err error) {
	var data []*process.WorkflowInstProc
	db := s.repository.GetDB().Model(&process.WorkflowInstProc{})
	// 流程的标题
	if req.Title != "" {
		db.Where("title like ?", "%"+req.Title+"%")
	}
	// 流程定义的名字
	if req.ProcDefName != "" {
		db.Where("proc_def_name like ?", "%"+req.ProcDefName+"%")
	}
	// 公司
	if req.Company != "" {
		db.Where("company like ?", "%"+req.Company+"%")
	}
	if req.ID != 0 {
		db.Where("id = ?", req.ID)
	}
	// 审核人
	if req.Candidate != "" {
		db.Where("candidate = ?", req.Candidate)
	}
	// 查询流程实例的列表
	err = db.Where("is_finished = 0").Count(&result.Total).Error
	if err != nil {
		return
	}
	err = db.Where("is_finished = 0").Offset(int(req.PageSize * req.Page)).Limit(int(req.PageSize)).Order("start_time desc").Find(&data).Error
	if err != nil {
		return
	}
	result.List = data
	result.PageSize = req.PageSize
	result.Page = req.Page
	return
}
