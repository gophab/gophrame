package i18n

import (
	"github.com/gophab/gophrame/core/context"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/i18n/config"
	"github.com/gophab/gophrame/core/starter"
)

func GetEnableLanguage() string {
	locale := context.GetContextValue("_LOCALE_")
	if locale == nil || locale.(string) == "" {
		return ""
	}
	return locale.(string)
}

func GetCurrentLanguage() string {
	locale := context.GetContextValue("_LOCALE_")
	if locale == nil || locale.(string) == "" {
		return defaultLanguage
	}
	return locale.(string)
}

func SetCurrentLanguage(locale string) {
	if locale == "" {
		context.RemoveContextValue("_LOCALE_")
	} else {
		context.SetContextValue("_LOCALE_", locale)
	}
}

func init() {
	starter.RegisterStarter(Start)
}

func Start() {
	if config.Setting.Enabled {
		if global.DB != nil {
			global.DB.Callback().Create().After("gorm:create").Register("LocaleUpdateHook", LocaleUpdateHook)
			global.DB.Callback().Update().After("gorm:update").Register("LocaleUpdateHook", LocaleUpdateHook)
			global.DB.Callback().Query().After("gorm:query").Register("LocaleLoadHook", LocaleLoadHook)
		}

		i18nManager = New()
		i18nManager.init()
	}
}
