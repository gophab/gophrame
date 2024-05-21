package json

import "encoding/json"

func String(obj interface{}) string {
	if bytes, err := json.Marshal(obj); err == nil {
		return string(bytes)
	}

	return ""
}

func Json(str string, obj interface{}) error {
	return json.Unmarshal([]byte(str), obj)
}
