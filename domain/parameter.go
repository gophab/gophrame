package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Parameter struct {
	Name         string `json:"name"`
	Type         string `json:"type,omitempty"`
	Description  string `json:"description,omitempty"`
	Value        string `json:"value,omitempty"`
	DefaultValue string `json:"defaultValue,omitempty"`
}

// Properties
type Parameters []*Parameter

// Value return json value, implement driver.Valuer interface
func (j Parameters) Value() (driver.Value, error) {
	if result, err := json.Marshal(j); err == nil {
		return string(result), nil
	} else {
		return nil, nil
	}
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *Parameters) Scan(value any) error {
	if value == nil {
		*j = Parameters(nil)
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

	return j.UnmarshalJSON(bytes)
}

// MarshalJSON to output non base64 encoded []byte
func (j Parameters) MarshalJSON() ([]byte, error) {
	var v []*Parameter = j
	return json.Marshal(v)
}

// UnmarshalJSON to deserialize []byte
func (j *Parameters) UnmarshalJSON(b []byte) error {
	var v = make([]*Parameter, 0)
	err := json.Unmarshal(b, &v)
	*j = v
	return err
}

func (j Parameters) String() string {
	var v []*Parameter = j
	if bytes, err := json.Marshal(v); err == nil {
		return string(bytes)
	} else {
		return ""
	}
}

// GormDataType gorm common data type
func (Parameters) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Parameters) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (js Parameters) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
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
func (j Parameters) Parameter(name string) *Parameter {
	for _, p := range j {
		if p.Name == name {
			return p
		}
	}
	return nil
}

type ParametersEnabled struct {
	Parameters *Parameters `gorm:"column:parameters;type:json" json:"parameters,omitempty"` /* 扩展信息 */
}

func (p *ParametersEnabled) Parameter(name string) *Parameter {
	if p.Parameters == nil {
		return &Parameter{Name: name}
	}
	return p.Parameters.Parameter(name)
}

func (p *ParametersEnabled) AddParameter(name string) *Parameter {
	var parameters []*Parameter
	if p.Parameters == nil {
		parameters = make([]*Parameter, 0)
	} else {
		parameters = *p.Parameters
	}
	result := &Parameter{
		Name: name,
	}
	parameters = append(parameters, result)
	*p.Parameters = parameters
	return result
}

func (p *ParametersEnabled) SetParameter(name string, defaultValue string) {
	if p.Parameters == nil {
		p.AddParameter(name).DefaultValue = defaultValue
	} else {
		parameter := p.Parameters.Parameter(name)
		if nil == parameter {
			p.AddParameter(name).DefaultValue = defaultValue
		} else {
			parameter.Value = defaultValue
			parameter.DefaultValue = defaultValue
		}
	}
}
