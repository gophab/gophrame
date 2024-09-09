package domain

import "github.com/gophab/gophrame/core/i18n"

type LocaleField struct {
	*i18n.LocaleFieldValue
}

func (*LocaleField) TableName() string {
	return "sys_locale_field"
}
