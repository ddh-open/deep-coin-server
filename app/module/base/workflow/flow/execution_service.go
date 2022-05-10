package flow

import (
	"devops-http/app/module/base/utils"
	"devops-http/app/module/workflow/model/identity"
	"devops-http/app/module/workflow/model/node"
	processModel "devops-http/app/module/workflow/model/process"
	taskModel "devops-http/app/module/workflow/model/task"
	identitySvc "devops-http/app/module/workflow/service/identity"
	taskSvc "devops-http/app/module/workflow/service/task"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// MoveStage MoveStage
// 流程流转
func MoveStage(nodeInfos []*node.WorkflowNodeInfo, link *identity.WorkflowIdentityLink, candidate string, pass bool, tx *gorm.DB) (err error) {
	err = identitySvc.SaveIdentityTx(link, tx)
	if err != nil {
		return err
	}
	if pass {
		link.Step = link.Step + 1
		if link.Step-1 > len(nodeInfos) {
			return errors.New("已经结束无法流转到下一个节点")
		}
	} else {
		link.Step = link.Step - 1
		if link.Step < 0 {
			return errors.New("处于开始位置，无法回退到上一个节点")
		}
	}
	// 指定下一步执行人
	if len(candidate) > 0 {
		nodeInfos[link.Step].Approved = candidate
	}
	// 判断下一流程： 如果是审批人是：抄送人
	// fmt.Println(nodeInfos[step].AproverType == flow.NodeTypes[flow.NOTIFIER])
	if nodeInfos[link.Step].ApprovedType == node.WorkflowNodeTypesList[node.NOTIFIER].Type {
		// 生成新的任务
		var task = &taskModel.WorkflowTask{
			NodeID:     node.WorkflowNodeTypesList[0].Type,
			Step:       link.Step,
			ProcInstID: link.ProcInstID,
			IsFinished: true,
		}
		task.IsFinished = true
		_, err = taskSvc.SaveTaskByTx(task, tx)
		if err != nil {
			return err
		}
		// 添加抄送人
		notifierLink := &identity.WorkflowIdentityLink{
			Group:      nodeInfos[link.Step].Approved,
			Type:       identity.KindsIdentity[identity.NOTIFIER],
			Step:       link.Step,
			ProcInstID: link.ProcInstID,
			Company:    link.Company,
			Comment:    "",
		}
		// 要判断抄送人， 是否已经存在了
		err = identitySvc.SaveIdentityTx(notifierLink, tx)
		if err != nil {
			return err
		}
		return MoveStage(nodeInfos, link, "", pass, tx)
	}
	if pass {
		return MoveToNextStage(nodeInfos, link, "", tx)
	}
	return MoveToPrevStage(nodeInfos, link, "", tx)
}

// MoveToNextStage MoveToNextStage
//通过
func MoveToNextStage(nodeInfos []*node.WorkflowNodeInfo, link *identity.WorkflowIdentityLink, candidate string, tx *gorm.DB) error {
	var currentTime = utils.FormatDate(time.Now())
	var task = getNewTask(nodeInfos, link.Step, link.ProcInstID) //新任务
	var processInst = &processModel.WorkflowInstProc{            // 流程实例要更新的字段
		NodeID:    nodeInfos[link.Step].NodeID,
		Candidate: nodeInfos[link.Step].Approved,
	}
	processInst.ID = link.ProcInstID
	if (link.Step + 1) != len(nodeInfos) { // 下一步不是【结束】
		// 生成新的任务
		taskId, err := taskSvc.SaveTaskByTx(task, tx)
		if err != nil {
			return err
		}
		// 添加candidate group
		notifierLink := &identity.WorkflowIdentityLink{
			Group:      nodeInfos[link.Step].Approved,
			Type:       node.WorkflowNodeTypesList[0].Type,
			Step:       link.Step,
			ProcInstID: processInst.ID,
			Company:    link.Company,
			TaskID:     taskId,
			Comment:    "",
		}
		err = identitySvc.SaveIdentityTx(notifierLink, tx)
		if err != nil {
			return err
		}
		// 更新流程实例
		processInst.TaskID = taskId
		err = UpdateProcInst(processInst, tx)
		if err != nil {
			return err
		}
	} else { // 最后一步直接结束
		// 生成新的任务
		task.IsFinished = true
		task.ClaimTime = currentTime
		taskId, err := taskSvc.SaveTaskByTx(task, tx)
		if err != nil {
			return err
		}
		// 删除候选用户组
		err = DelCandidateByProcInstID(link.ProcInstID, tx)
		if err != nil {
			return err
		}
		// 更新流程实例
		processInst.TaskID = taskId
		processInst.EndTime = currentTime
		processInst.IsFinished = true
		processInst.Candidate = "审批结束"
		err = UpdateProcInst(processInst, tx)
		if err != nil {
			return err
		}
	}
	return nil
}

// MoveToPrevStage MoveToPrevStage
// 驳回
func MoveToPrevStage(nodeInfos []*node.WorkflowNodeInfo, link *identity.WorkflowIdentityLink, candidate string, tx *gorm.DB) error {
	// 生成新的任务
	var task = getNewTask(nodeInfos, link.Step, link.ProcInstID) //新任务
	taskId, err := taskSvc.SaveTaskByTx(task, tx)
	if err != nil {
		return err
	}

	var processInst = &processModel.WorkflowInstProc{ // 流程实例要更新的字段
		NodeID:    nodeInfos[link.Step].NodeID,
		Candidate: nodeInfos[link.Step].Approved,
		TaskID:    taskId,
	}

	processInst.ID = link.ProcInstID

	err = UpdateProcInst(processInst, tx)

	if err != nil {
		return err
	}
	if link.Step == 0 { // 流程回到起始位置，注意起始位置为0,
		err = AddCandidateUserTx(nodeInfos[link.Step].Approved, link.Company, link.Step, taskId, link.ProcInstID, tx)
		if err != nil {
			return err
		}
		return nil
	}
	// 添加candidate group
	err = AddCandidateGroupTx(nodeInfos[link.Step].Approved, link.Company, link.Step, taskId, link.ProcInstID, tx)
	if err != nil {
		return err
	}
	return nil
}

// AddCandidateUserTx AddCandidateUserTx
// 添加候选用户
func AddCandidateUserTx(userID, company string, step int, taskID, procInstID uint, tx *gorm.DB) error {
	err := DelCandidateByProcInstID(procInstID, tx)
	if err != nil {
		return err
	}
	i := &identity.WorkflowIdentityLink{
		UserID:     userID,
		Type:       identity.KindsIdentity[0],
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return identitySvc.SaveIdentityTx(i, tx)
}

// AddCandidateGroupTx AddCandidateGroupTx
// 添加候选用户组
func AddCandidateGroupTx(group, company string, step int, taskID, procInstID uint, tx *gorm.DB) error {
	err := DelCandidateByProcInstID(procInstID, tx)
	if err != nil {
		return err
	}
	i := &identity.WorkflowIdentityLink{
		Group:      group,
		Type:       identity.KindsIdentity[0],
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return identitySvc.SaveIdentityTx(i, tx)
}

func getNewTask(nodeInfos []*node.WorkflowNodeInfo, step int, procInstID uint) *taskModel.WorkflowTask {
	var task = &taskModel.WorkflowTask{ // 新任务
		NodeID:        nodeInfos[step].NodeID,
		Step:          step,
		ProcInstID:    procInstID,
		MemberCount:   nodeInfos[step].MemberCount,
		UnCompleteNum: nodeInfos[step].MemberCount,
		ActType:       nodeInfos[step].ActType,
	}
	return task
}

// DelCandidateByProcInstID DelCandidateByProcInstID
// 删除历史候选人
func DelCandidateByProcInstID(procInstID uint, tx *gorm.DB) error {
	return tx.Where("proc_inst_id=? and type=?", procInstID, identity.KindsIdentity[0]).Delete(&identity.WorkflowIdentityLink{}).Error
}

// ExistsByProcInstIDAndGroup 判断是否已经存在WorkflowIdentityLink
func ExistsByProcInstIDAndGroup(procInstID int, group string, kind int, tx gorm.DB) (bool, error) {
	var count int64
	err := tx.Model(&identity.WorkflowIdentityLink{}).Where("identitylink.proc_inst_id=? and identitylink.group=? and identitylink.type=?", procInstID, group, identity.KindsIdentity[kind]).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func UpdateProcInst(workflowInstProcData *processModel.WorkflowInstProc, tx *gorm.DB) error {
	return tx.Model(&processModel.WorkflowInstProc{}).Where("id = ?", workflowInstProcData.ID).Updates(&workflowInstProcData).Error
}
