package container

import (
	"log"
	"strings"
	"sync"

	"github.com/wjshen/gophrame/errors"
)

// 定义一个全局键值对存储容器
var sMap sync.Map

// CreateContainersFactory 创建一个容器工厂
func CreateContainersFactory() *container {
	return &container{}
}

// 定义一个容器结构体
type container struct {
}

// Set  1.以键值对的形式将代码注册到容器
func (c *container) Set(key string, value interface{}) (res bool) {
	if _, exists := c.KeyIsExists(key); !exists {
		sMap.Store(key, value)
		res = true
	} else {
		log.Fatal(errors.ERROR_CONTAINER_KEY_ALREADY_EXISTS, errors.GetErrorMessage(errors.ERROR_CONTAINER_KEY_ALREADY_EXISTS), ",请解决键名重复问题,相关键："+key)
	}
	return
}

// Delete  2.删除
func (c *container) Delete(key string) {
	sMap.Delete(key)
}

// Get 3.传递键，从容器获取值
func (c *container) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

// KeyIsExists 4. 判断键是否被注册
func (c *container) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *container) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPre) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
