package config

import (
	"time"
	"unsafe"

	"github.com/wjshen/gophrame/core/logger"

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
		if value, _ := json.MarshalToString(config); value != "" {
			logger.Debug("Load application configuration: ", value)
		}

		// Second to json
		for key, value := range configs {
			if key == "ROOT" {
				// "ROOT" node
				logger.Debug("Load system configuration")
				err = UnmarshalFromNode(config, value.Setting)
			} else if node, ok := config[key]; ok {
				// Setting node
				logger.Debug("Load module configuration: ", key)
				err = UnmarshalFromNode(node, value.Setting)
			}

			if err != nil {
				logger.Error(err.Error())
				break
			}
		}
	}

	return err
}
