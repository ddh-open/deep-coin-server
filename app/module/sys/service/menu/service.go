package menu

import (
	"devops-http/app/contract"
	"devops-http/app/module/base"
	"devops-http/app/module/base/request"
	"devops-http/app/module/base/response"
	"devops-http/app/module/base/sys"
	"devops-http/app/module/sys/model/menu"
	"devops-http/app/module/sys/model/operation"
	"devops-http/app/module/sys/model/role"
	"devops-http/framework"
	contract2 "devops-http/framework/contract"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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
	db.AutoMigrate(&menu.DevopsSysMenu{}, &menu.DevopsSysMenuParameter{}, &menu.DevopsSysMenuBtn{}, &operation.DevopsSysOperationRecord{})
	return &Service{base.NewRepository(db)}
}

func (s *Service) GetRepository() *base.Repository {
	return s.repository
}

func (s *Service) SetRepository(model interface{}) *base.Repository {
	return s.repository.SetRepository(model)
}

func (s *Service) getMenuTreeMap(roleId string) (err error, treeMap map[string][]menu.DevopsSysMenu) {
	var allMenus []menu.DevopsSysMenu
	var btns []menu.DevopsSysMenuBtn
	treeMap = make(map[string][]menu.DevopsSysMenu)
	err = s.repository.GetDB().Where("authority_id = ?", roleId).Order("sort").Preload("Parameters").Find(&allMenus).Error
	if err != nil {
		return
	}
	err = s.repository.GetDB().Where("authority_id = ?", roleId).Preload("SysBaseMenuBtn").Find(&btns).Error
	if err != nil {
		return
	}
	var btnMap = make(map[uint]map[string]string)
	for _, v := range btns {
		if btnMap[v.DevopsSysMenuID] == nil {
			btnMap[v.DevopsSysMenuID] = make(map[string]string)
		}
		btnMap[v.DevopsSysMenuID][v.Name] = roleId
	}
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

func (s *Service) GetList(req request.SearchMenusParams) (list interface{}, err error) {
	var result response.PageResult
	result.Page = req.Page
	result.PageSize = req.PageSize

	limit := int(result.PageSize)
	offset := int(result.PageSize * (result.Page - 1))
	db := s.repository.GetDB().Model(&menu.DevopsSysMenu{})
	var menuList []menu.DevopsSysMenu
	db = db.Where("parent_id =  0")
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}

	if req.Component != "" {
		db = db.Where("component = ?", req.Component)
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
			orderMap["path"] = true
			orderMap["name"] = true
			orderMap["title"] = true
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
			err = db.Order(OrderStr).Find(&menuList).Error
		} else {
			err = db.Order("id").Find(&menuList).Error
		}
	}
	for i := 0; i < len(menuList); i++ {
		err = s.getBaseChildrenList(&menuList[i])
	}
	result.List = menuList
	return result, err
}

func (s *Service) getChildrenList(menuData *menu.DevopsSysMenu, treeMap map[string][]menu.DevopsSysMenu) (err error) {
	menuData.Children = treeMap[strconv.Itoa(int(menuData.ID))]
	for i := 0; i < len(menuData.Children); i++ {
		err = s.getChildrenList(&menuData.Children[i], treeMap)
	}
	return err
}

func (s *Service) getBaseChildrenList(menuData *menu.DevopsSysMenu) (err error) {
	var children []menu.DevopsSysMenu
	s.repository.SetRepository(&menu.DevopsSysMenu{}).GetDB().Where("parent_id = ?", menuData.ID).Find(&children)
	menuData.Children = children
	for i := 0; i < len(menuData.Children); i++ {
		err = s.getBaseChildrenList(&menuData.Children[i])
	}
	return err
}

func (s *Service) AddBaseMenu(menuData menu.DevopsSysMenuEntity) error {
	if !errors.Is(s.repository.SetRepository(&menu.DevopsSysMenu{}).GetDB().Where("name = ?", menuData.Name).First(&menu.DevopsSysMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err := s.repository.GetDB().Create(&menuData.DevopsSysMenu).Error
	if err != nil {
		return err
	}
	return err
}

func (s *Service) getBaseMenuTreeMap() (err error, treeMap map[string][]menu.DevopsSysMenu) {
	var allMenus []menu.DevopsSysMenu
	treeMap = make(map[string][]menu.DevopsSysMenu)
	err = s.repository.GetDB().Order("sort").Preload("MenuBtn").Preload("Parameters").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

func (s *Service) GetBaseMenuTree() (err error, menus []menu.DevopsSysMenu) {
	err, treeMap := s.getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = s.getChildrenList(&menus[i], treeMap)
	}
	return err, menus
}

func (s *Service) AddMenuToRole(userToken *base.TokenUser, req sys.RelativeRoleMenuRequest, cabin contract.Cabin) (err error) {
	if len(req.MenuIds) <= 0 {
		return errors.New("菜单id为空！")
	}
	cabin.GetCabin().ClearPolicy()
	_, err = cabin.GetCabin().RemoveFilteredNamedPolicy("p", 0, req.RoleId, userToken.CurrentDomain, "", "MENUS")
	if err != nil {
		return err
	}
	var rule [][]string
	for i := range req.MenuIds {
		rule = append(rule, []string{req.RoleId, userToken.CurrentDomain, req.MenuIds[i], "MENUS"})
	}
	// 给角色添加菜单树
	cabin.GetCabin().ClearPolicy()
	_, err = cabin.GetCabin().AddPolicies(rule)
	return err
}

// GetMenuByRole 根据角色获取菜单
func (s *Service) GetMenuByRole(userToken *base.TokenUser, id string, cabin contract.Cabin) (err error, menus []string) {
	var roleData role.DevopsSysRole
	s.repository.SetRepository(&role.DevopsSysRole{}).GetDB().First(&roleData, "id = ?", id)
	if roleData.ID <= 0 {
		err = errors.New("未找到该角色！")
		return
	}
	cabin.GetCabin().ClearPolicy()
	data := cabin.GetCabin().GetFilteredNamedPolicy("p", 0, cast.ToString(roleData.ID), userToken.CurrentDomain, "", "MENUS")
	for _, datum := range data {
		menus = append(menus, datum[2])
	}
	return err, menus
}

// GetMenuByUser 根据用户获取菜单
func (s *Service) GetMenuByUser(tokenUser *base.TokenUser, cabin contract.Cabin) (err error, menus []menu.DevopsSysMenu) {
	list, err := cabin.GetCabin().GetImplicitResourcesForUser(tokenUser.Uuid, tokenUser.CurrentDomain)
	var menusIdStr []string
	for _, str := range list {
		if str[3] == "MENUS" {
			menusIdStr = append(menusIdStr, str[2])
		}
	}
	var allMenus []menu.DevopsSysMenu
	err = s.repository.SetRepository(&menu.DevopsSysMenu{}).GetDB().Order("sort").Preload("MenuBtn").Preload("Parameters").Preload("Parameters").Find(&allMenus, menusIdStr).Error
	if err != nil {
		return
	}
	treeMap := make(map[string][]menu.DevopsSysMenu, 0)
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = s.getChildrenList(&menus[i], treeMap)
	}
	return err, menus
}

func (s *Service) DeleteBaseMenu(id sys.DeleteById) (err error) {
	err = s.repository.GetDB().Preload("MenuBtn").Preload("Parameters").Where("parent_id in (?)", id).First(&menu.DevopsSysMenu{}).Error
	if err != nil {
		var menuData []menu.DevopsSysMenu
		err = s.repository.GetDB().Where("id in (?)", id.Ids).Find(&menuData).Error
		if err != nil {
			return err
		}
		for _, v := range menuData {
			err = s.repository.GetDB().Delete(&menu.DevopsSysMenu{}, "id = ?", v.ID).Error
			err = s.repository.GetDB().Delete(&menu.DevopsSysMenuParameter{}, "devops_sys_menu_id = ?", v.ID).Error
			err = s.repository.GetDB().Delete(&menu.DevopsSysMenuBtn{}, "devops_sys_menu_id = ?", v.ID).Error
			// 删除相关的权限 此处预留
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("此菜单存在子菜单不可删除")
	}
	return err
}

// UpdateBaseMenu 更新路由菜单
func (s *Service) UpdateBaseMenu(menuData menu.DevopsSysMenuEntity) (err error) {
	var oldMenu menu.DevopsSysMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menuData.KeepAlive
	upDateMap["close_tab"] = menuData.CloseTab
	upDateMap["default_menu"] = menuData.DefaultMenu
	upDateMap["parent_id"] = menuData.ParentId
	upDateMap["path"] = menuData.Path
	upDateMap["name"] = menuData.Name
	upDateMap["hidden"] = menuData.Hidden
	upDateMap["component"] = menuData.Component
	upDateMap["title"] = menuData.Title
	upDateMap["icon"] = menuData.Icon
	upDateMap["sort"] = menuData.Sort

	err = s.repository.GetDB().Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", menuData.ID).Find(&oldMenu)
		if oldMenu.Name != menuData.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menuData.ID, menuData.Name).First(&menu.DevopsSysMenu{}).Error, gorm.ErrRecordNotFound) {
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := tx.Unscoped().Delete(&menu.DevopsSysMenuParameter{}, "devops_sys_menu_id = ?", menuData.ID).Error
		if txErr != nil {
			return txErr
		}
		txErr = tx.Unscoped().Delete(&menu.DevopsSysMenuBtn{}, "devops_sys_menu_id = ?", menuData.ID).Error
		if txErr != nil {
			return txErr
		}
		if len(menuData.Parameters) > 0 {
			for k := range menuData.Parameters {
				menuData.Parameters[k].DevopsSysMenuID = menuData.ID
			}
			txErr = tx.Create(&menuData.Parameters).Error
			if txErr != nil {
				return txErr
			}
		}

		if len(menuData.MenuBtn) > 0 {
			for k := range menuData.MenuBtn {
				menuData.MenuBtn[k].DevopsSysMenuID = menuData.ID
			}
			txErr = tx.Create(&menuData.MenuBtn).Error
			if txErr != nil {
				return txErr
			}
		}

		txErr = db.Updates(upDateMap).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return err
}

// GetBaseMenuById 根据id 查询菜单
func (s *Service) GetBaseMenuById(id string) (err error, menuData menu.DevopsSysMenu) {
	err = s.repository.GetDB().Preload("MenuBtn").Preload("Parameters").Where("id = ?", id).First(&menuData).Error
	if err != nil {
		return
	}
	err = s.getBaseChildrenList(&menuData)
	return
}
