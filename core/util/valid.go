package util

import (
	"regexp"

	"github.com/astaxie/beego/validation"
)

// @Summary   更新用户
// @Tags  users
// @Accept json
// @Produce  json
// @Param   body  body   models.User   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {string} json
// @Router /api/v1/users/:id  [PUT]
var internationalPhonePattern = regexp.MustCompile(`^(([\+0])?\d{1,4}(\-)?)?\d{5,13}$`)

// Tel check telephone struct
type InternationalTelephone struct {
	validation.Match
	Key string
}

func NewInternationalTelephoneValidator(field string) validation.Validator {
	return InternationalTelephone{validation.Match{Regexp: internationalPhonePattern}, field}
}
