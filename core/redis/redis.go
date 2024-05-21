package redis

import "github.com/wjshen/gophrame/core/redis/config"

func init() {
	initRedisClientPool(config.Setting.Database)
}
