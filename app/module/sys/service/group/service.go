package group

import (
	"context"
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/group"
	"devops-http/app/module/sys/model/menu"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"devops-http/resources/proto/userGrpc"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Service struct {
	repository *base.Repository
}

func NewService(c framework.Container) *Service {
	db, err := c.MustMake(contract2.ORMKey).(contract2.ORMService).GetDB()
	logger := c.MustMake(contract2.LogKey).(contract2.Log)
	if err != nil {
		logger.Error("service 获取db出错： err", zap.Error(err))
	}
	db.AutoMigrate(&group.DevopsSysGroup{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetGroupById(id string) (menuData *group.DevopsSysGroup, err error) {
	err = s.repository.GetDB().Where("id = ?", id).First(menuData).Error
	if err != nil {
		return
	}
	err = s.getBaseChildrenList(menuData)
	return
}

func (s *Service) getChildrenList(menuData *group.DevopsSysGroup, treeMap map[string][]group.DevopsSysGroup) (err error) {
	menuData.Children = treeMap[strconv.Itoa(int(menuData.ID))]
	for i := 0; i < len(menuData.Children); i++ {
		err = s.getChildrenList(&menuData.Children[i], treeMap)
	}
	return err
}

func (s *Service) getBaseChildrenList(groupData *group.DevopsSysGroup) (err error) {
	var children []group.DevopsSysGroup
	s.repository.SetRepository(&menu.DevopsSysMenu{}).GetDB().Where("parent_id = ?", groupData.ID).Find(&children)
	groupData.Children = children
	for i := 0; i < len(groupData.Children); i++ {
		err = s.getBaseChildrenList(&groupData.Children[i])
	}
	return err
}

func (s *Service) GetGroupList(req request.SearchGroupParams) (result response.PageResult, err error) {
	result.Page = req.Page
	result.PageSize = req.PageSize
	limit := int(result.PageSize)
	offset := int(result.PageSize * (result.Page - 1))
	db := s.repository.GetDB().Model(&group.DevopsSysGroup{})
	var groupList []group.DevopsSysGroup
	db = db.Where("parent_id =  0")
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Linkman != "" {
		db = db.Where("linkman LIKE ?", "%"+req.Linkman+"%")
	}

	err = db.Count(&result.Total).Error

	if err != nil {
		return
	} else {
		db = db.Limit(limit).Offset(offset)
		if req.OrderKey != "" {
			var OrderStr string
			// 设置有效排序key 防止sql注入
			// 感谢 Tom4t0 提交漏洞信息
			orderMap := make(map[string]bool, 4)
			orderMap["id"] = true
			orderMap["name"] = true
			if orderMap[req.OrderKey] {
				if req.Desc {
					OrderStr = req.OrderKey + " desc"
				} else {
					OrderStr = req.OrderKey
				}
			} else { // didn't matched any order key in `orderMap`
				err = fmt.Errorf("非法的排序字段: %v", req.OrderKey)
				return
			}
			err = db.Order(OrderStr).Find(&groupList).Error
		} else {
			err = db.Order("id").Find(&groupList).Error
		}
	}
	for i := 0; i < len(groupList); i++ {
		err = s.getBaseChildrenList(&groupList[i])
	}
	result.List = groupList
	return result, err
}

func (s *Service) AddGroup(req group.DevopsSysGroup) error {
	if !errors.Is(s.repository.SetRepository(&group.DevopsSysGroup{}).GetDB().Where("name = ?", req.Name).First(&group.DevopsSysGroup{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err := s.repository.GetDB().Create(&req).Error
	if err != nil {
		return err
	}
	return err
}

func (s *Service) ModifyGroup(req group.DevopsSysGroup) (err error) {
	var oldMenu group.DevopsSysGroup
	upDateMap := make(map[string]interface{})
	upDateMap["linkman"] = req.Linkman
	upDateMap["linkman_no"] = req.LinkmanNo
	upDateMap["alias"] = req.Alias
	upDateMap["remark"] = req.Remark
	upDateMap["name"] = req.Name
	upDateMap["parent_id"] = req.ParentID
	upDateMap["sort"] = req.Sort
	upDateMap["enable"] = req.Enable

	err = s.repository.GetDB().Model(&group.DevopsSysGroup{}).Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", req.ID).Find(&oldMenu)
		if oldMenu.Name != req.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", req.ID, req.Name).First(&menu.DevopsSysMenu{}).Error, gorm.ErrRecordNotFound) {
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return err
}

func (s *Service) DeleteGroup(ids string, grpcService contract.ServiceGrpc, param ...interface{}) error {
	var idsInt []int64
	if strings.Contains(ids, ",") {
		for _, s2 := range strings.Split(ids, ",") {
			idsInt = append(idsInt, cast.ToInt64(s2))
		}
	} else {
		idsInt = append(idsInt, cast.ToInt64(ids))
	}
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	client := userGrpc.NewServiceGroupClient(conn)
	resp, err := client.GroupDelete(context.Background(), &userGrpc.IdsRequest{
		Ids: idsInt,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}

func (s *Service) AddResourcesToGroup(request []request.CabinInReceive, grpcService contract.ServiceGrpc, param ...interface{}) error {
	conn, err := grpcService.GetGrpc("grpc.user")
	if err != nil {
		return err
	}
	defer conn.Close()
	requestBytes, err := json.Marshal(&request)
	if err != nil {
		return err
	}
	client := userGrpc.NewServiceCabinClient(conn)
	resp, err := client.CabinRuleAdd(context.Background(), &userGrpc.BytesRequest{
		Data: requestBytes,
	})
	if err != nil {
		return err
	}
	if resp.GetCode() != 200 {
		err = errors.Wrap(err, resp.GetMsg())
	}
	return err
}
