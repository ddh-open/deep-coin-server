package casbin

import (
	"github.com/casbin/casbin/v2"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/ddh-open/gin/framework"
	"github.com/ddh-open/gin/framework/contract"
	"gorm.io/gorm"
)

func NewCaBinService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	config := c.MustMake(contract.ConfigKey).(contract.Config)
	db, err := c.MustMake(contract.ORMKey).(contract.ORMService).GetDB()
	if err != nil {
		return nil, err
	}
	a, err := gormAdapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	syncedEnforcer, err := casbin.NewSyncedEnforcer(config.GetString("app.casbin.model-path"), a)
	if err != nil {
		return nil, err
	}
	_ = syncedEnforcer.LoadPolicy()
	return &ServiceCabin{
		Enforcer: syncedEnforcer,
		db:       db,
	}, nil
}

type ServiceCabin struct {
	Enforcer *casbin.SyncedEnforcer
	db       *gorm.DB
}

// GetCabin @function: Cabin
//@description: 持久化到数据库  引入自定义规则
//@return: *Cabin.Enforcer
func (sc *ServiceCabin) GetCabin() *casbin.SyncedEnforcer {
	_ = sc.Enforcer.LoadPolicy()
	return sc.Enforcer
}

func (sc *ServiceCabin) ClearCabin(v int, p ...string) bool {
	e := sc.GetCabin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

func (sc *ServiceCabin) UpdateCabinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := sc.db.Model(gormAdapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}
