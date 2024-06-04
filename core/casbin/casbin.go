package casbin

import (
	"errors"
	"net/http"

	"github.com/gophab/gophrame/core/casbin/config"
	"github.com/gophab/gophrame/core/database"
	"github.com/gophab/gophrame/core/webservice/response"
	"github.com/gophab/gophrame/service"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

const (
	// casbin 初始化可能的错误
	ErrorCasbinCanNotUseDbPtr         string = "casbin 的初始化基于gorm 初始化后的数据库连接指针，程序检测到 gorm 连接指针无效，请检查数据库配置！"
	ErrorCasbinCreateAdaptFail        string = "casbin NewAdapterByDBUseTableName 发生错误："
	ErrorCasbinCreateEnforcerFail     string = "casbin NewEnforcer 发生错误："
	ErrorCasbinNewModelFromStringFail string = "NewModelFromString 调用时出错："
	ErrorsCasbinNoAuthorization       string = "Casbin 鉴权未通过，请在后台检查 casbin 设置参数"
)

// 创建 Casbin Enforcer(执行器):
// 初始化：
// 1. 自动创建数据库表
// 2. 加载策略
// 3. 设置自动加载策略间隔
func InitCasbinEnforcer() (*casbin.SyncedEnforcer, error) {
	var Enforcer *casbin.SyncedEnforcer

	adapter, err := gormAdapter.NewAdapterByDBUseTableName(database.DB(), config.Setting.TablePrefix, config.Setting.TableName)
	if err != nil {
		return nil, errors.New(ErrorCasbinCreateAdaptFail)
	}

	if model, err := model.NewModelFromString(config.Setting.ModelConfig); err != nil {
		return nil, errors.New(ErrorCasbinNewModelFromStringFail + err.Error())
	} else {
		if Enforcer, err = casbin.NewSyncedEnforcer(model, adapter); err != nil {
			return nil, errors.New(ErrorCasbinCreateEnforcerFail)
		}
		_ = Enforcer.LoadPolicy()
		if config.Setting.AutoLoadPolicyInterval > 0 {
			Enforcer.StartAutoLoadPolicy(config.Setting.AutoLoadPolicyInterval)
		}

		return Enforcer, nil
	}
}

// casbin 鉴权失败，返回 405 方法不允许访问
func ErrorCasbinAuthFail(c *gin.Context, msg interface{}) {
	response.ErrorMessage(c, http.StatusForbidden, http.StatusMethodNotAllowed, ErrorsCasbinNoAuthorization)
}

// CasbinService负责更新Casbin Enforce数据至
type CasbinService struct {
	Enforcer *casbin.SyncedEnforcer `inject:"enforcer"`
}

func (s *CasbinService) onUserEvent() {

	service.GetUserService().LoadPolicy(res.Id)

}
