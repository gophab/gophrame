package aliyun

import (
	"github.com/gophab/gophrame/core"
	"github.com/gophab/gophrame/core/http"
	"github.com/gophab/gophrame/core/identify"
	"github.com/gophab/gophrame/core/identify/aliyun/config"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/domain"
	"github.com/gophab/gophrame/errors"

	"gorm.io/gorm"
)

type RealNameHistory struct {
	domain.Logable
	Name       string `gorm:"column:name;primaryKey"`
	Mobile     string `gorm:"column:mobile;primaryKey"`
	NationalId string `gorm:"column:national_id;primaryKey"`
	Result     bool   `gorm:"column:result;default:false"`
}

func (*RealNameHistory) TableName() string {
	return "t_real_name_history"
}

type VerifyService struct {
	Db      *gorm.DB `inject:"database"`
	history map[string]bool
}

var verifyService = &VerifyService{
	history: make(map[string]bool),
}

type Result struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    *struct {
		Result  string `json:"result"`
		OrderNo string `json:"orderNo"`
		Desc    string `json:"desc"`
	} `json:"data"`
}

func init() {
	inject.InjectValue("verifyService", verifyService)
}

func Start() {
	if config.Setting.Enabled {
		identify.Api = verifyService
	}
}

func (s *VerifyService) RealNameVerify(name string, mobile string, nationalId string) (bool, error) {
	if nationalId != "" {
		if v, b := s.history[name+":"+mobile+":"+nationalId]; b {
			return v, nil
		}
	}
	if v, b := s.history[name+":"+mobile]; b {
		if !v {
			return v, nil
		}
	}

	var r = false
	if nationalId != "" {
		if res := s.Db.Model(&RealNameHistory{}).Where("name = ?", name).Where("mobile = ?", mobile).Where("national_id = ?", nationalId).Select("result").First(&r); res.Error == nil && res.RowsAffected > 0 {
			s.history[name+":"+mobile+":"+nationalId] = r
			return r, nil
		}
	}

	if res := s.Db.Model(&RealNameHistory{}).Where("name = ?", name).Where("mobile = ?", mobile).Where("national_id = ?", "*").Select("result").First(&r); res.Error == nil && res.RowsAffected > 0 {
		if !r { //
			s.history[name+":"+mobile] = false
			return r, nil
		}
		s.history[name+":"+mobile] = true
		if nationalId == "" {
			return true, nil
		}
	}

	var req *http.HttpRequest = http.NewHttpRequest(config.Setting.Base).
		HEADER("Authorization", "APPCODE "+config.Setting.AppCode).
		HEADER("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8").
		PROXY(config.Setting.Proxy)

	if nationalId == "" {
		// 二要素认证
		if config.Setting.TwoFactorUrl == "" {
			req.POST(config.Setting.Url)
		} else {
			req.POST(config.Setting.TwoFactorUrl)
		}
		req.BODY(core.M{
			"mobile": mobile,
			"name":   name,
		})
	} else {
		// 三要素认证
		if config.Setting.ThreeFactorUrl == "" {
			req.POST(config.Setting.Url)
		} else {
			req.POST(config.Setting.ThreeFactorUrl)
		}
		req.BODY(core.M{
			"mobile": mobile,
			"name":   name,
			"idcard": nationalId,
		})
	}

	var result Result
	if status, err := req.DoForm().ResultTo(&result); err == nil {
		logger.Debug("Identify api result: ", json.String(result))

		if status == 200 {
			if result.Success && result.Data != nil && result.Data.Result == "0" {
				s.history[name+":"+mobile] = true
				s.Db.Save(&RealNameHistory{
					Name:       name,
					Mobile:     mobile,
					NationalId: "*",
					Result:     true,
				})
				if nationalId != "" {
					s.Db.Save(&RealNameHistory{
						Name:       name,
						Mobile:     mobile,
						NationalId: nationalId,
						Result:     true,
					})
					s.history[name+":"+mobile+":"+nationalId] = true
				}
				return true, nil
			} else {
				if nationalId == "" {
					s.history[name+":"+mobile] = false
					s.Db.Save(&RealNameHistory{
						Name:       name,
						Mobile:     mobile,
						NationalId: "*",
						Result:     false,
					})
				} else {
					s.history[name+":"+mobile+":"+nationalId] = false
					s.Db.Save(&RealNameHistory{
						Name:       name,
						Mobile:     mobile,
						NationalId: nationalId,
						Result:     true,
					})
				}
				return false, nil
			}
		} else {
			logger.Error("Identify api error: ", "API result error - ", status)
			return false, errors.New(status, "API result error")
		}
	} else {
		logger.Error("Identify api error: ", err.Error())
		return false, err
	}
}
