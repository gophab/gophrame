package json

import "encoding/json"

func String(obj any) string {
	if obj == nil {
		return ""
	}

	if bytes, err := json.Marshal(obj); err == nil {
		return string(bytes)
	}

	return ""
}

func Json(str string, obj any) error {
	return json.Unmarshal([]byte(str), obj)
}
