package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Property struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Scope        string `json:"scope"` /* 作用范围：USER、ADMIN、ALL */
	Value        any    `json:"value"`
	DefaultValue any    `json:"defaultValue"`
}

func (p *Property) GetValue() any {
	if p.Value != "" {
		return p.Value
	}

	return p.DefaultValue
}

func (p *Property) Default(v any) *Property {
	p.DefaultValue = v
	return p
}

func (p *Property) AsBool() bool {
	return asBool(p.GetValue())
}

func (p *Property) AsInt() int {
	return asInt(p.GetValue())
}

func (p *Property) AsString() string {
	return asString(p.GetValue())
}

func asBool(v any) bool {
	if v == nil {
		return false
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Int:
		return v.(int) != 0
	case reflect.Int64:
		return v.(int64) != 0
	case reflect.Float64:
		return v.(float64) != 0
	case reflect.Float32:
		return v.(float32) != 0
	case reflect.String:
		return strings.ToLower(v.(string)) == "true"
	case reflect.Bool:
		return v.(bool)
	default:
		return true
	}
}

func asInt(v any) int {
	if v == nil {
		return 0
	}

	switch v := v.(type) {
	case int:
		return v
	case string:
		r, _ := strconv.Atoi(v)
		return r
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case float32, float64:
		return int(math.Round(v.(float64)))
	default:
		return 0
	}
}

func asString(v any) string {
	if v == nil {
		return ""
	}

	switch v := v.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case float32, float64:
		return strconv.FormatFloat(v.(float64), 'f', 4, 64)
	default:
		bytes, _ := json.Marshal(v)
		if bytes != nil {
			return string(bytes)
		} else {
			return ""
		}
	}
}

// Properties
type Properties map[string]any

// Value return json value, implement driver.Valuer interface
func (j Properties) Value() (driver.Value, error) {
	if result, err := json.Marshal(j); err == nil {
		return string(result), nil
	} else {
		return nil, nil
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *Properties) Scan(value any) error {
	if value == nil {
		*j = Properties(nil)
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result = make(map[string]any)
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (j Properties) MarshalJSON() ([]byte, error) {
	var v map[string]any = j
	return json.Marshal(v)
}

// UnmarshalJSON to deserialize []byte
func (j *Properties) UnmarshalJSON(b []byte) error {
	var v = make(map[string]any)
	err := json.Unmarshal(b, &v)
	*j = v
	return err
}

func (j Properties) String() string {
	var v map[string]any = j
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	} else {
		return ""
	}
}

// GormDataType gorm common data type
func (Properties) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Properties) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (js Properties) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(js) == 0 {
		return gorm.Expr("NULL")
	}

	data, _ := js.MarshalJSON()

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

// Value return json value, implement driver.Valuer interface
func (j Properties) Property(name string) *Property {
	if v, b := j[name]; b {
		return &Property{
			Name:  name,
			Value: v,
		}
	}
	return &Property{
		Name: name,
	}
}

type PropertiesEnabled struct {
	Properties *Properties `gorm:"column:properties;type:json" json:"properties,omitempty" i18n:"yes"` /* 扩展信息 */
}

func (p *PropertiesEnabled) Property(name string) *Property {
	if p.Properties == nil {
		return &Property{DefaultValue: nil}
	}
	return p.Properties.Property(name)
}

func (p *PropertiesEnabled) SetProperty(name string, value any) {
	var properties map[string]any
	if p.Properties == nil {
		properties = make(map[string]any)
		p.Properties = &Properties{}
		*p.Properties = properties
	} else {
		properties = *p.Properties
	}
	properties[name] = value
}
