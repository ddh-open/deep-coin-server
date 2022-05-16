package group

import (
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/sys/model/group"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
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
	err = db.AutoMigrate(&base.DevopsSysGroup{}, &group.DevopsSysGroupRelativeUser{})
	// 建立多对多的关系表
	err = db.SetupJoinTable(&base.DevopsSysGroup{}, "Users", &group.DevopsSysGroupRelativeUser{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetGroupById(id string) (menuData *base.DevopsSysGroup, err error) {
	err = s.repository.GetDB().Where("id = ?", id).First(menuData).Error
	if err != nil {
		return
	}
	err = s.getBaseChildrenList(menuData)
	return
}

func (s *Service) getChildrenList(groupData *base.DevopsSysGroup, treeMap map[string][]base.DevopsSysGroup) (err error) {
	var users []base.DevopsSysUser
	s.repository.GetDB().Model(groupData).Association("Users").Find(&users)
	groupData.Children = treeMap[strconv.Itoa(int(groupData.ID))]
	for i := 0; i < len(groupData.Children); i++ {
		err = s.getChildrenList(&groupData.Children[i], treeMap)
	}
	return err
}

func (s *Service) getBaseChildrenList(groupData *base.DevopsSysGroup) (err error) {
	var children []base.DevopsSysGroup
	var users []base.DevopsSysUser
	s.repository.SetRepository(&base.DevopsSysGroup{}).GetDB().Where("parent_id = ?", groupData.ID).Find(&children)
	groupData.Children = children
	s.repository.GetDB().Model(groupData).Association("Users").Find(&users)
	groupData.Users = users
	for i := 0; i < len(groupData.Children); i++ {
		err = s.getBaseChildrenList(&groupData.Children[i])
	}
	return err
}

func (s *Service) getWithOutUserChildrenList(groupData *base.DevopsSysGroup) (err error) {
	var children []base.DevopsSysGroup
	s.repository.SetRepository(&base.DevopsSysGroup{}).GetDB().Where("parent_id = ?", groupData.ID).Find(&children)
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
	db := s.repository.GetDB().Model(&base.DevopsSysGroup{})
	var groupList []base.DevopsSysGroup
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
	for i, _ := range groupList {
		var users []base.DevopsSysUser
		s.repository.GetDB().Model(&groupList[i]).Association("Users").Find(&users)
		groupList[i].Users = users
		if groupList[i].ParentID == 0 {
			groupList[i].ParentName = "根节点"
			continue
		}
		s.repository.GetDB().Model(&base.DevopsSysGroup{}).Where("id = ?", groupList[i].ParentID).Select("name").Find(&groupList[i].ParentName)
	}
	result.List = groupList
	return result, err
}

func (s *Service) GetGroupTree() (result response.PageResult, err error) {
	db := s.repository.GetDB().Model(&base.DevopsSysGroup{})
	var groupList []base.DevopsSysGroup
	err = db.Count(&result.Total).Error
	db.Where("parent_id = ?", 0)
	if err != nil {
		return
	} else {
		err = db.Order("sort").Find(&groupList).Error
	}
	for i := 0; i < len(groupList); i++ {
		err = s.getWithOutUserChildrenList(&groupList[i])
	}
	result.List = groupList
	return result, err
}

func (s *Service) AddGroup(req base.DevopsSysGroup) error {
	if !errors.Is(s.repository.SetRepository(&base.DevopsSysGroup{}).GetDB().Where("name = ?", req.Name).First(&base.DevopsSysGroup{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err := s.repository.GetDB().Create(&req).Error
	if err != nil {
		return err
	}
	return err
}

func (s *Service) ModifyGroup(req base.DevopsSysGroup) (err error) {
	var oldMenu base.DevopsSysGroup
	upDateMap := make(map[string]interface{})
	upDateMap["linkman"] = req.Linkman
	upDateMap["linkman_no"] = req.LinkmanNo
	upDateMap["alias"] = req.Alias
	upDateMap["remark"] = req.Remark
	upDateMap["name"] = req.Name
	upDateMap["parent_id"] = req.ParentID
	upDateMap["sort"] = req.Sort
	upDateMap["enable"] = req.Enable

	err = s.repository.GetDB().Model(&base.DevopsSysGroup{}).Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", req.ID).Find(&oldMenu)
		if oldMenu.Name != req.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", req.ID, req.Name).First(&base.DevopsSysGroup{}).Error, gorm.ErrRecordNotFound) {
				return errors.New("存在相同name修改失败")
			}
		}
		err = db.Updates(upDateMap).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Service) DeleteGroup(req request.ReqById) (err error) {
	err = s.repository.GetDB().Where("parent_id in (?)", req.Ids).First(&base.DevopsSysGroup{}).Error
	if err != nil {
		var groupData []base.DevopsSysGroup
		err = s.repository.GetDB().Where("id in (?)", req.Ids).Find(&groupData).Error
		if err != nil {
			return err
		}
		for _, v := range groupData {
			err = s.repository.GetDB().Model(&v).Association("Users").Clear()
			// 删除相关的权限 此处预留
			if err != nil {
				return err
			}
			err = s.repository.GetDB().Unscoped().Delete(&base.DevopsSysGroup{}, "id = ?", v.ID).Error
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("此菜单存在子菜单不可删除")
	}
	return
}

func (s *Service) AddUserToGroup(req request.GroupRelativeUserRequest) (err error) {
	var groupData base.DevopsSysGroup
	err = s.repository.GetDB().Model(&base.DevopsSysGroup{}).Where("id = ?", req.GroupId).First(&groupData).Error
	if err != nil {
		return
	}
	var users []base.DevopsSysUser
	err = s.repository.GetDB().Model(&base.DevopsSysUser{}).Where("id in (?)", req.UserIds).Find(&users).Error
	if err != nil {
		return
	}
	return s.repository.GetDB().Model(&groupData).Association("Users").Append(users)
}

func (s *Service) DeleteUserToGroup(req request.GroupRelativeUserRequest) (err error) {
	var groupData base.DevopsSysGroup
	err = s.repository.GetDB().Model(&base.DevopsSysGroup{}).Where("id = ?", req.GroupId).First(&groupData).Error
	if err != nil {
		return
	}
	var users []base.DevopsSysUser
	err = s.repository.GetDB().Model(&base.DevopsSysUser{}).Where("id in (?)", req.UserIds).Find(&users).Error
	if err != nil {
		return
	}
	return s.repository.GetDB().Model(&groupData).Association("Users").Delete(users)
}
