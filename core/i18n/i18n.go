package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gophab/gophrame/core/context"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/starter"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type LocaleFieldValue struct {
	EntityName string `gorm:"column:entity_name;primaryKey" json:"entityName"`
	EntityId   string `gorm:"column:entity_id;primaryKey" json:"entityId"`
	Name       string `gorm:"column:name;primaryKey" json:"name"`
	Locale     string `gorm:"column:locale;default:en;primaryKey" json:"locale"`
	Value      string `gorm:"column:value" json:"value"`
}

type I18nEnabled struct {
	LocaleFields []*LocaleFieldValue `gorm:"-" json:"-"`
}

type Translator interface {
	StoreTranslations(translations []*LocaleFieldValue)
	LoadTranslations(locale, entityName string, entityId ...string) map[string][]*LocaleFieldValue
}

type I18nFactory struct {
	Translator Translator `inject:"translator"`
}

var i18nFactory = &I18nFactory{}

func init() {
	inject.InjectValue("i18nFactory", i18nFactory)
	starter.RegisterStarter(Start)
}

func Start() {
	if global.DB != nil {
		global.DB.Callback().Create().After("gorm:create").Register("LocaleUpdateHook", LocaleUpdateHook)
		global.DB.Callback().Update().After("gorm:update").Register("LocaleUpdateHook", LocaleUpdateHook)
		global.DB.Callback().Query().After("gorm:query").Register("LocaleLoadHook", LocaleLoadHook)
	}
}

func getTagSection(tag, key string) string {
	segs := strings.Split(tag, ";")
	for i := 0; i < len(segs); i++ {
		if strings.HasPrefix(segs[i], key+":") {
			return segs[i]
		}
	}
	return ""
}

func buildLocaleField(db *gorm.DB, item reflect.Value, field *schema.Field, locale string, columns []string) []*LocaleFieldValue {
	ctx := db.Statement.Context
	if v, isZero := field.ValueOf(ctx, item); !isZero {
		// 构造
		idField := db.Statement.Schema.LookUpField("Id")
		id, _ := idField.ValueOf(ctx, item)
		if fmt.Sprint(id) != "" {
			vt := reflect.ValueOf(v).Kind()
			switch vt {
			case reflect.Array, reflect.Slice:
				if bs, err := json.Marshal(v); err == nil {
					return []*LocaleFieldValue{
						{
							EntityName: db.Statement.Schema.ModelType.Name(),
							EntityId:   fmt.Sprint(id),
							Name:       field.Name,
							Locale:     locale,
							Value:      string(bs),
						},
					}
				}
			case reflect.Map, reflect.Struct: //
				if len(columns) > 0 {
					// 只保存字段
					var results = make([]*LocaleFieldValue, 0)
					for _, column := range columns {
						results = append(results, &LocaleFieldValue{
							EntityName: db.Statement.Schema.ModelType.Name(),
							EntityId:   fmt.Sprint(id),
							Name:       field.Name + "." + column,
							Locale:     locale,
							Value:      fmt.Sprint(v),
						})
					}
					return results
				} else {
					if bs, err := json.Marshal(v); err == nil {
						return []*LocaleFieldValue{
							{
								EntityName: db.Statement.Schema.ModelType.Name(),
								EntityId:   fmt.Sprint(id),
								Name:       field.Name,
								Locale:     locale,
								Value:      string(bs),
							},
						}
					}
				}
			default:
				return []*LocaleFieldValue{
					{
						EntityName: db.Statement.Schema.ModelType.Name(),
						EntityId:   fmt.Sprint(id),
						Name:       field.Name,
						Locale:     locale,
						Value:      fmt.Sprint(v),
					},
				}
			}
		}
	}
	return nil
}

// 1. store locale fields
func LocaleUpdateHook(db *gorm.DB) {
	if i18nFactory.Translator == nil {
		return
	}

	locale := context.GetContextValue("_LOCALE_")
	if locale == nil || locale.(string) == "" {
		return
	}

	if db.Statement.Schema == nil {
		return
	}

	if len(db.Statement.Schema.Fields) == 0 {
		return
	}

	var localeFields = make([]*LocaleFieldValue, 0)

	for _, field := range db.Statement.Schema.Fields {
		if tag, b := field.Tag.Lookup("i18n"); b {
			// 1. field 是基本数据类型
			properties := strings.Split(getTagSection(tag, "property"), ",")
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					item := db.Statement.ReflectValue.Index(i)
					fieldValues := buildLocaleField(db, item, field, locale.(string), properties)
					if fieldValues != nil {
						localeFields = append(localeFields, fieldValues...)
					}
				}
			case reflect.Struct:
				fieldValues := buildLocaleField(db, db.Statement.ReflectValue, field, locale.(string), properties)
				if fieldValues != nil {
					localeFields = append(localeFields, fieldValues...)
				}
			}
		}
	}

	if len(localeFields) > 0 {
		// translator.StoreTranslation()
		i18nFactory.Translator.StoreTranslations(localeFields)
	}
}

// 2. get locale fields
func LocaleLoadHook(db *gorm.DB) {
	if i18nFactory.Translator == nil {
		return
	}

	locale := context.GetContextValue("_LOCALE_")
	if locale == nil || locale.(string) == "" {
		return
	}

	if db.Statement.Schema == nil {
		return
	}

	if len(db.Statement.Schema.Fields) == 0 {
		return
	}

	var localeFields = make([]*schema.Field, 0)
	for _, field := range db.Statement.Schema.Fields {
		if v, b := field.Tag.Lookup("i18n"); b && v == "yes" {
			localeFields = append(localeFields, field)
		}
	}

	ctx := db.Statement.Context
	if len(localeFields) > 0 {
		var ids []string = make([]string, 0)
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				item := db.Statement.ReflectValue.Index(i)
				if idField := db.Statement.Schema.LookUpField("Id"); idField != nil {
					if id, b := idField.ValueOf(ctx, item); !b {
						ids = append(ids, fmt.Sprint(id))
					}
				}
			}
		case reflect.Struct:
			if idField := db.Statement.Schema.LookUpField("Id"); idField != nil {
				if id, b := idField.ValueOf(ctx, db.Statement.ReflectValue); !b {
					ids = append(ids, fmt.Sprint(id))
				}
			}
		}

		if len(ids) == 0 {
			return
		}

		// translator.StoreTranslation()
		localeFieldValues := i18nFactory.Translator.LoadTranslations(
			locale.(string),
			db.Statement.Schema.ModelType.Name(),
			ids...,
		)

		if len(localeFieldValues) > 0 {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					item := db.Statement.ReflectValue.Index(i)
					idField := db.Statement.Schema.LookUpField("Id")
					if id, b := idField.ValueOf(ctx, item); !b {
						filtered := localeFieldValues[fmt.Sprint(id)]
						if len(filtered) > 0 {
							for _, fieldValue := range filtered {
								if field := db.Statement.Schema.LookUpField(fieldValue.Name); field != nil {
									field.Set(ctx, item, fieldValue.Value)
								}
							}
						}
					}
				}
			case reflect.Struct:
				for _, fieldValue := range localeFieldValues[ids[0]] {
					if field := db.Statement.Schema.LookUpField(fieldValue.Name); field != nil {
						field.Set(ctx, db.Statement.ReflectValue, fieldValue.Value)
					}
				}
			}
		}
	}
}
