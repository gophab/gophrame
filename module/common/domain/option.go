package domain

import "strings"

type Option struct {
	Name        string `gorm:"column:name" json:"name"`
	Value       string `gorm:"column:value" json:"value"`
	ValueType   string `gorm:"column:value_type" json:"valueType"`
	Description string `gorm:"column:description;default:null" json:"description"`
}

func OptionsToMap(sysEnvs []Option) map[string]interface{} {
	var result = make(map[string]interface{})
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

func getMapObject(m *map[string]interface{}, key string) *map[string]interface{} {
	var segs = strings.Split(key, ".")
	var result *map[string]interface{} = m
	for _, s := range segs {
		if v, ok := (*result)[s]; ok {
			if vm, ok := v.(map[string]interface{}); ok {
				result = &vm
			} else {
				var r = make(map[string]interface{})
				(*result)[s] = r
				result = &r
			}
		} else {
			var r = make(map[string]interface{})
			(*result)[s] = r
			result = &r
		}
	}
	return result
}
