package node

import (
	"container/list"
	"devops-http/app/module/base"
)

var WorkflowNodeTypesList []WorkflowNodeType
var WorkflowNodeKindsList []WorkflowNodeType

const (
	START = iota
	APPROVE
	ROUTE
	EXEC
	NOTIFIER
	CONDITION
	END
)

const (
	// RANGE 条件类型: 范围
	RANGE = iota
	// VALUE 条件类型： 值
	VALUE
)

// ActionConditionTypes 所有条件类型
var ActionConditionTypes = [...]string{RANGE: "range_condition", VALUE: "value_condition"}

type WorkflowNodeParams struct {
	base.DevopsModel
	Key    string `json:"name,omitempty"`
	Source int    `json:"source,omitempty"`
	Type   string `json:"type,omitempty"`
	Remark string `json:"remark,omitempty"`
}

// Node represents a specific logical unit of processing and routing
// in a workflow.
// 流程中的一个节点
type Node struct {
	Name           string      `json:"name,omitempty"`
	Type           string      `json:"type,omitempty"`
	NodeID         string      `json:"nodeId,omitempty"`
	PrevID         string      `json:"prevId,omitempty"`
	ChildNode      *Node       `json:"childNode,omitempty"`
	ConditionNodes []*Node     `json:"conditionNodes,omitempty"`
	Properties     *Properties `json:"properties,omitempty"`
}

type Properties struct {
	// ONE_BY_ONE 代表依次审批
	ActivateType     string             `json:"activateType,omitempty"`
	AgreeAll         bool               `json:"agreeAll,omitempty"`
	Conditions       [][]*ConditionNode `json:"conditions,omitempty"`
	ActionRules      []*ActionRule      `json:"actionRules,omitempty"`
	NoneActionAction string             `json:"noneActionAction,omitempty"`
	Script           string             `json:"script,omitempty"`
}

type ConditionNode struct {
	Type       string `json:"type,omitempty"`
	ParamKey   string `json:"paramKey,omitempty"`
	ParamLabel string `json:"paramLabel,omitempty"`
	IsEmpty    bool   `json:"isEmpty,omitempty"`
	// 类型为range
	LowerBound      string `json:"lowerBound,omitempty"`
	LowerBoundEqual string `json:"lowerBoundEqual,omitempty"`
	UpperBoundEqual string `json:"upperBoundEqual,omitempty"`
	UpperBound      string `json:"upperBound,omitempty"`
	BoundEqual      string `json:"boundEqual,omitempty"`
	Unit            string `json:"unit,omitempty"`
	// 类型为 value
	ParamValues []string    `json:"paramValues,omitempty"`
	OriValue    []string    `json:"oriValue,omitempty"`
	Conditions  []*CondNode `json:"conditions,omitempty"`
}

type CondNode struct {
	Type  string    `json:"type,omitempty"`
	Value string    `json:"value,omitempty"`
	Attrs *UserNode `json:"attrs,omitempty"`
}

type UserNode struct {
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// WorkflowNodeInfo 节点信息
type WorkflowNodeInfo struct {
	NodeID       string `json:"nodeId"`
	Type         string `json:"type"`
	Approved     string `json:"approved"`
	ApprovedType string `json:"approvedType"`
	MemberCount  int8   `json:"memberCount"`
	Level        int8   `json:"level"`
	ActType      string `json:"actType"`
}

type ActionRule struct {
	Type       string `json:"type,omitempty"`
	LabelNames string `json:"labelNames,omitempty"`
	Labels     int    `json:"labels,omitempty"`
	IsEmpty    bool   `json:"isEmpty,omitempty"`
	// 表示需要通过的人数 如果是会签
	MemberCount int8 `json:"memberCount,omitempty"`
	// and 表示会签 or表示或签，默认为或签
	ActType string `json:"actType,omitempty"`
	Level   int8   `json:"level,omitempty"`
	AutoUp  bool   `json:"autoUp,omitempty"`
}

type WorkflowNodeType struct {
	base.DevopsModel
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Kind         int    `json:"kind,omitempty"`
	RelativeKind int    `json:"relativeKind,omitempty"`
	Remark       string `json:"remark,omitempty"`
}

func (n *Node) AddToExecutionList(list *list.List) {
	switch n.Type {
	case "approver", "notifier":
		approve := n.Properties.ActionRules[0].LabelNames
		list.PushBack(WorkflowNodeInfo{
			NodeID:       n.NodeID,
			Type:         n.Properties.ActivateType,
			Approved:     approve,
			ApprovedType: n.Type,
			MemberCount:  n.Properties.ActionRules[0].MemberCount,
			ActType:      n.Properties.ActionRules[0].ActType,
		})
	case "exec":
		approve := n.Properties.ActionRules[0].LabelNames
		list.PushBack(WorkflowNodeInfo{
			NodeID:       n.NodeID,
			Type:         n.Properties.ActivateType,
			Approved:     approve,
			ApprovedType: n.Type,
			MemberCount:  n.Properties.ActionRules[0].MemberCount,
			ActType:      n.Properties.ActionRules[0].ActType,
		})
		break
	default:
	}
}
