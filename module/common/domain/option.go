package domain

import "strings"

type Option struct {
	Name        string `gorm:"column:name;primaryKey" json:"name"`
	Value       string `gorm:"column:value" json:"value"`
	ValueType   string `gorm:"column:value_type" json:"valueType"`
	Description string `gorm:"column:description;default:null" json:"description"`
}

func OptionsToMap(sysEnvs []Option) map[string]any {
	var result = make(map[string]any)
	for _, se := range sysEnvs {
		var pos = strings.LastIndex(se.Name, ".")
		if pos >= 0 {
			var r = getMapObject(&result, se.Name[:(pos-1)])
			(*r)[se.Name[(pos+1):]] = se.Value
		} else {
			result[se.Name] = se.Value
		}
	}

	return result
}

func getMapObject(m *map[string]any, key string) *map[string]any {
	var segs = strings.Split(key, ".")
	var result *map[string]any = m
	for _, s := range segs {
		if v, ok := (*result)[s]; ok {
			if vm, ok := v.(map[string]any); ok {
				result = &vm
			} else {
				var r = make(map[string]any)
				(*result)[s] = r
				result = &r
			}
		} else {
			var r = make(map[string]any)
			(*result)[s] = r
			result = &r
		}
	}
	return result
}
