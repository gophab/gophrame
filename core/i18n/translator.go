package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/domain"
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

type I18nPropertiesEnabled interface {
	I18nProperties() string
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

func buildLocaleField(db *gorm.DB, item reflect.Value, field *schema.Field, locale string) []*LocaleFieldValue {
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
			case reflect.Pointer:
				return []*LocaleFieldValue{
					{
						EntityName: db.Statement.Schema.ModelType.Name(),
						EntityId:   fmt.Sprint(id),
						Name:       field.Name,
						Locale:     locale,
						Value:      fmt.Sprintf("%v", reflect.ValueOf(v).Elem()),
					},
				}
			default:
				return []*LocaleFieldValue{
					{
						EntityName: db.Statement.Schema.ModelType.Name(),
						EntityId:   fmt.Sprint(id),
						Name:       field.Name,
						Locale:     locale,
						Value:      fmt.Sprintf("%v", v),
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

	locale := GetEnableLanguage()
	if locale == "" {
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
		if _, b := field.Tag.Lookup("i18n"); b {
			// 1. field 是基本数据类型
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					item := db.Statement.ReflectValue.Index(i)
					fieldValues := buildLocaleField(db, item, field, locale)
					if fieldValues != nil {
						localeFields = append(localeFields, fieldValues...)
					}
				}
			case reflect.Struct:
				fieldValues := buildLocaleField(db, db.Statement.ReflectValue, field, locale)
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

func getFields(str string) (results []string) {
	if strings.TrimSpace(str) == "" {
		return nil
	}

	results = strings.Split(strings.TrimSpace(str), ",")
	return
}

func getPropertyFields(str string) (results map[string][]string) {
	results = make(map[string][]string)
	if str == "" {
		return
	}
	var segs = strings.Split(str, ";")
	for _, seg := range segs {
		kv := strings.Split(strings.TrimSpace(seg), ":")
		if len(kv) > 1 {
			var key = kv[0]
			var fields = getFields(kv[1])
			results[strings.ToLower(key)] = fields
		}
	}
	return
}

// data[field] = value
func setSchemaLocaleField(ctx context.Context, data reflect.Value, field *schema.Field, value any) {
	if field.DataType == "json" {
		// 属性字段：type == domain.Properties
		var properties []string
		if propertiesEnabled, ok := data.Interface().(I18nPropertiesEnabled); ok {
			properties = getPropertyFields(propertiesEnabled.I18nProperties())[strings.ToLower(field.Name)]
		} else if tag, b := field.Tag.Lookup("i18n"); b {
			properties = getFields(getTagSection(tag, "properties"))
		}

		if len(properties) > 0 {
			// 1. 获取原始数据
			var oldJsonValue = make(map[string]any)
			var jsonValue = make(map[string]any)
			var err error = nil
			if ov, b := field.ValueOf(ctx, data); !b && ov != nil {
				oldJsonValue = *ov.(*domain.Properties)
			}

			if v, b := value.(string); b && v != "" {
				err = json.Unmarshal([]byte(v), &jsonValue)
			}

			if err == nil {
				var updated = false
				for _, p := range properties {
					if nv, b := jsonValue[p]; b && nv != nil {
						oldJsonValue[p] = nv
						updated = true
					}
				}

				if updated {
					if bytes, err := json.Marshal(oldJsonValue); err == nil {
						field.Set(ctx, data, string(bytes))
						// logger.Debug("Update i18n properties: ", strings.Join(properties, ","))
					}
				}
			}

			return
		}

		// 2. Struct 字段
		// if field.Name == "Dimensions" {
		// 	logger.Debug("Dimensions")
		// }
	}
	field.Set(ctx, data, value)
}

// field = value
func setLocaleField(data reflect.Value, fieldName string, value any) {
	if dataField := data.FieldByName(fieldName); dataField.IsValid() && !dataField.IsZero() {
		if dataField.Type().Name() == "json" {
			var properties []string
			if propertiesEnabled, ok := data.Interface().(I18nPropertiesEnabled); ok {
				properties = getPropertyFields(propertiesEnabled.I18nProperties())[strings.ToLower(fieldName)]
			} else {
				if field, b := data.Type().FieldByName(fieldName); b {
					if tag, b := field.Tag.Lookup("i18n"); b {
						properties = getFields(getTagSection(tag, "properties"))
					}
				}
			}

			if len(properties) > 0 {
				// 1. 获取原始数据
				var oldJsonValue = make(map[string]any)
				var jsonValue = make(map[string]any)
				var err error = nil
				if ov := dataField.String(); ov != "" {
					err = json.Unmarshal([]byte(ov), &oldJsonValue)
				}

				if err == nil {
					if v, b := value.(string); b && v != "" {
						err = json.Unmarshal([]byte(v), &jsonValue)
					}
				}

				if err == nil {
					var updated = false
					for _, p := range properties {
						if nv, b := jsonValue[p]; b && nv != nil {
							oldJsonValue[p] = jsonValue[p]
							updated = true
						}
					}

					if updated {
						if bytes, err := json.Marshal(oldJsonValue); err == nil {
							dataField.Set(reflect.ValueOf(string(bytes)))
							// logger.Debug("Update i18n properties: ", strings.Join(properties, ","))
						}
					}

				}
				return
			}
		}

		dataField.Set(reflect.ValueOf(value))
	}
}

// 2. get locale fields
func LocaleLoadHook(db *gorm.DB) {
	if i18nFactory.Translator == nil {
		return
	}

	locale := GetEnableLanguage()
	if locale == "" {
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
	idField := db.Statement.Schema.LookUpField("Id")
	if len(localeFields) > 0 {
		var ids []string = make([]string, 0)
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				item := db.Statement.ReflectValue.Index(i)
				var v = reflect.Indirect(item)
				if v.Kind() == reflect.Struct {
					if v.Type() == db.Statement.Schema.ModelType {
						if idField != nil {
							if id, b := idField.ValueOf(ctx, item); !b {
								ids = append(ids, fmt.Sprint(id))
							}
						}
					} else if id := v.FieldByName("Id"); id.IsValid() && !id.IsZero() {
						ids = append(ids, fmt.Sprint(id))
					}
				}
			}
		case reflect.Struct:
			var v = reflect.Indirect(db.Statement.ReflectValue)
			if v.Type() == db.Statement.Schema.ModelType {
				if idField != nil {
					if id, b := idField.ValueOf(ctx, db.Statement.ReflectValue); !b {
						ids = append(ids, fmt.Sprint(id))
					}
				}
			} else if id := v.FieldByName("Id"); id.IsValid() && !id.IsZero() {
				ids = append(ids, fmt.Sprint(id))
			}
		}

		if len(ids) == 0 {
			return
		}

		// translator.StoreTranslation()
		localeFieldValues := i18nFactory.Translator.LoadTranslations(
			locale,
			db.Statement.Schema.ModelType.Name(),
			ids...,
		)

		if len(localeFieldValues) > 0 {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					item := db.Statement.ReflectValue.Index(i)
					var v = reflect.Indirect(item)
					if v.Kind() == reflect.Struct {
						if v.Type() == db.Statement.Schema.ModelType {
							if idField != nil {
								if id, b := idField.ValueOf(ctx, item); !b {
									filtered := localeFieldValues[fmt.Sprint(id)]
									if len(filtered) > 0 {
										for _, fieldValue := range filtered {
											if field := db.Statement.Schema.LookUpField(fieldValue.Name); field != nil {
												setSchemaLocaleField(ctx, item, field, fieldValue.Value)
											}
										}
									}
								}
							}
						} else if id := v.FieldByName("Id"); id.IsValid() && !id.IsZero() {
							filtered := localeFieldValues[fmt.Sprint(id.Interface())]
							if len(filtered) > 0 {
								for _, fieldValue := range filtered {
									setLocaleField(v, fieldValue.Name, fieldValue.Value)
								}
							}
						}
					}
				}
			case reflect.Struct:
				var v = reflect.Indirect(db.Statement.ReflectValue)
				if v.Type() == db.Statement.Schema.ModelType {
					for _, fieldValue := range localeFieldValues[ids[0]] {
						if field := db.Statement.Schema.LookUpField(fieldValue.Name); field != nil {
							setSchemaLocaleField(ctx, db.Statement.ReflectValue, field, fieldValue.Value)
						}
					}
				} else {
					for _, fieldValue := range localeFieldValues[ids[0]] {
						setLocaleField(v, fieldValue.Name, fieldValue.Value)
					}
				}
			}
		}
	}
}
