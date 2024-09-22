package i18n

import (
	"github.com/gophab/gophrame/core/context"
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
