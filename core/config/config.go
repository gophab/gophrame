package config

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gophab/gophrame/core/logger"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

// Module Config Settings
type ConfigSetting struct {
	Name        string
	Description string
	Setting     interface{}
}

var configs = make(map[string]ConfigSetting)

func RegisterConfig(name string, setting interface{}, desc string) {
	configs[name] = ConfigSetting{
		Name:        name,
		Description: desc,
		Setting:     setting,
	}
}

var configChangeCallbacks = make([]func(), 0)

func RegisterConfigChangeCallback(f func()) {
	configChangeCallbacks = append(configChangeCallbacks, f)
}

type JsonExtension struct {
	jsoniter.DummyExtension
}

type DurationDecoder struct{}

func (*DurationDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	switch iter.WhatIsNext() {
	case jsoniter.StringValue:
		if d, err := time.ParseDuration(iter.ReadString()); err == nil {
			*(*time.Duration)(ptr) = time.Duration(d)
		}
	case jsoniter.NilValue:
		*((*time.Duration)(ptr)) = time.Duration(0)
	default:
		*((*time.Duration)(ptr)) = time.Duration(iter.ReadInt64())
	}
}

func (*JsonExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	if typ.AssignableTo(reflect2.TypeOf(time.Duration(0))) {
		return &DurationDecoder{}
	}
	return nil
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	json.RegisterExtension(&JsonExtension{})
	RegisterConfigChangeCallback(func() {
		loadConfig()
	})
}

func InitConfig() error {
	return loadConfig()
}

func UnmarshalFromNode(node any, out interface{}) error {
	if data, err := json.Marshal(node); err == nil {
		// Use json to
		return json.Unmarshal(data, out)
	} else {
		return err
	}
}

func loadConfig() error {
	var config = make(map[string]interface{})

	// First load into map[]
	var err error
	if err = InitYamlConfig(&config); err == nil {
		// Logger
		if text, _ := json.MarshalToString(config); text != "" {
			logger.Debug("Load application configuration: ", text)
		}

		// Second to json
		for key, value := range configs {
			if key == "ROOT" {
				// "ROOT" node
				err = UnmarshalFromNode(config, value.Setting)
			} else if node, ok := getConfigNode(config, key); ok {
				// Setting node
				logger.Debug("Load module configuration: ", key)
				err = UnmarshalFromNode(node, value.Setting)
			}

			if err != nil {
				logger.Error("Load configuration error: ", key, err.Error())
				break
			} else {
				if text, _ := json.MarshalToString(value.Setting); text != "" {
					logger.Debug("Load configuration: ", key, text)
				}
			}
		}
	}

	return err
}

func getConfigNode(config interface{}, path string) (interface{}, bool) {
	if config == nil {
		return nil, false
	}

	segs := strings.SplitN(path, ".", 1)
	if len(segs) == 2 {
		if node, b := getConfigNode(config, segs[0]); b {
			return getConfigNode(node, segs[1])
		}
	} else {
		switch reflect.TypeOf(config).Kind() {
		case reflect.Map:
			value := reflect.ValueOf(config).MapIndex(reflect.ValueOf(path))
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return value.Interface(), true
			}
		case reflect.Array:
			if index, err := strconv.ParseInt(path, 10, 32); err == nil {
				array := reflect.ValueOf(config)
				if 0 <= index && index <= int64(array.Len()) {
					value := reflect.ValueOf(config).Index(int(index))
					if value.IsValid() && !value.IsNil() && !value.IsZero() {
						return value.Interface(), true
					}
				}
			}
		case reflect.Struct:
			value := reflect.ValueOf(config)
			field := value.FieldByName(path)
			if field.IsValid() && !field.IsNil() && !field.IsZero() {
				return field.Interface(), true
			}
		case reflect.Ptr, reflect.Interface:
			value := reflect.ValueOf(config).Elem()
			if value.IsValid() && !value.IsNil() && !value.IsZero() {
				return getConfigNode(value.Interface(), path)
			}
		default:
		}
	}
	return nil, false
}
