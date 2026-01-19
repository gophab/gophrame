package identify

import (
	"errors"

	"github.com/gophab/gophrame/core/global"
)

type IdentifyApi interface {
	RealNameVerify(name string, mobile string, nationalId string) (bool, error)
}

var Api IdentifyApi

func TwoFactorIdentify(name, mobile string) (bool, error) {
	if name == "" || mobile == "" {
		return false, errors.New("要素不足")
	}

	if global.Debug {
		return true, nil
	}

	if Api != nil {
		return Api.RealNameVerify(name, mobile, "")
	} else {
		return false, nil
	}
}

func ThreeFactorIdentify(name, mobile, id string) (bool, error) {
	if mobile == "" {
		return TwoFactorIdentify(name, id)
	}

	if name == "" || id == "" {
		return false, errors.New("要素不足")
	}

	if global.Debug {
		return true, nil
	}

	if Api != nil {
		return Api.RealNameVerify(name, mobile, id)
	} else {
		return false, nil
	}
}

func RealNameIdentify(name, mobile, nationalId string) (bool, error) {
	if nationalId == "" {
		return TwoFactorIdentify(name, mobile)
	}

	if name == "" || nationalId == "" {
		return false, errors.New("要素不足")
	}

	if global.Debug {
		return true, nil
	}

	if Api != nil {
		return Api.RealNameVerify(name, mobile, nationalId)
	} else {
		return false, nil
	}
}

func init() {

}
