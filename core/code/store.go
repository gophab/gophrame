package code

import (
	"sync"
	"time"

	"github.com/gophab/gophrame/core/code/config"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/redis"
	"github.com/gophab/gophrame/errors"

	"github.com/patrickmn/go-cache"
)

type CodeStore interface {
	CreateRequest(id string) error
	CreateCode(id string, scene string, code string) error
	GetCode(id string, scene string, remove bool) (string, bool)
	RemoveCode(id string, scene string)
}

type CacheCodeStore struct {
	CodeStore
	codeCache *cache.Cache
	reqCache  *cache.Cache
}

func CreateCacheCodeStore(config *config.CodeStoreSetting) (*CacheCodeStore, error) {
	result := &CacheCodeStore{
		codeCache: cache.New(config.ExpireIn, time.Minute),
		reqCache:  cache.New(config.RequestInterval, time.Second*15),
	}
	return result, nil
}

func (s *CacheCodeStore) CreateRequest(id string) error {
	return s.reqCache.Add(id, time.Now(), 0)
}

func (s *CacheCodeStore) CreateCode(id string, scene string, code string) error {
	s.codeCache.Set(id+":"+scene, code, 0)
	return nil
}

func (s *CacheCodeStore) GetCode(id string, scene string, remove bool) (string, bool) {
	value, ok := s.codeCache.Get(id + ":" + scene)
	if ok {
		if remove {
			s.RemoveCode(id, scene)
		}
		return value.(string), true
	}
	return "", false
}

func (s *CacheCodeStore) RemoveCode(id string, scene string) {
	s.codeCache.Delete(id + ":" + scene)
}

type ExpireCode struct {
	Code       string
	Expiration time.Time
}

type MemoryCodeStore struct {
	data            sync.Map
	requestInterval time.Duration
	expireIn        time.Duration
}

func CreateMemoryCodeStore(config *config.CodeStoreSetting) (*MemoryCodeStore, error) {
	result := &MemoryCodeStore{data: sync.Map{}, requestInterval: config.RequestInterval, expireIn: config.ExpireIn}

	// 清除过期的验证码
	go func() {
		for {
			result.data.Range(func(key, value interface{}) bool {
				if value.(ExpireCode).Expiration.Before(time.Now()) {
					result.data.Delete(key)
				}
				return true
			})
			time.Sleep(time.Second * 60)
		}
	}()

	return result, nil
}

func (s *MemoryCodeStore) CreateRequest(id string) error {
	value, b := s.data.Load("req:" + id)
	if b && value.(ExpireCode).Expiration.After(time.Now()) {
		return errors.New("调用过于频繁")
	}

	s.data.Store("req:"+id, ExpireCode{
		Code:       "1",
		Expiration: time.Now().Add(time.Second * 60), /* 三分钟过期 */
	})
	return nil
}

func (s *MemoryCodeStore) CreateCode(id string, scene string, code string) error {
	s.data.Store("code:"+id+":"+scene, ExpireCode{
		Code:       code,
		Expiration: time.Now().Add(time.Second * 180), /* 三分钟过期 */
	})
	return nil
}

func (s *MemoryCodeStore) GetCode(id string, scene string, remove bool) (string, bool) {
	value, ok := s.data.Load("code:" + id + ":" + scene)
	if ok {
		if remove {
			s.RemoveCode(id, scene)
		}
		if time.Now().Before(value.(ExpireCode).Expiration) {
			return value.(ExpireCode).Code, true
		}
	}
	return "", false
}

func (s *MemoryCodeStore) RemoveCode(phone string, scene string) {
	s.data.Delete("code:" + phone + ":" + scene)
}

const (
	CODE_REDIS_KEY_PREFIX    = "code:phone:"
	REQUEST_REDIS_KEY_PREFIX = "request:phone:"
)

type RedisCodeStore struct {
	redisClient     *redis.RedisClient
	keyPrefix       string
	requestInterval time.Duration
	expireIn        time.Duration
}

func CreateRedisCodeStore(config *config.CodeStoreSetting) (result CodeStore, err error) {
	logger.Debug("Using redis code store")
	result = &RedisCodeStore{
		redisClient:     redis.GetOneRedisClientIndex(config.Redis.Database),
		keyPrefix:       config.Redis.KeyPrefix,
		requestInterval: config.RequestInterval,
		expireIn:        config.ExpireIn,
	}
	return result, nil
}

func (s *RedisCodeStore) CreateRequest(id string) error {
	if b, err := s.redisClient.Execute("EXISTS", s.RedisKey(REQUEST_REDIS_KEY_PREFIX, id, "req")); err != nil {
		return err
	} else if b.(int64) != 0 {
		// 键值存在
		return errors.New("请求频繁")
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(REQUEST_REDIS_KEY_PREFIX, id, "req"), s.requestInterval.Seconds(), "1"); err != nil {
		return err
	}

	return nil
}

func (s *RedisCodeStore) CreateCode(id string, scene string, code string) error {
	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(CODE_REDIS_KEY_PREFIX, id, scene), s.expireIn.Seconds(), code); err != nil {
		return err
	}
	return nil
}

func (s *RedisCodeStore) GetCode(id string, scene string, remove bool) (string, bool) {
	if code, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(CODE_REDIS_KEY_PREFIX, id, scene))); err != nil {
		return "", false
	} else {
		if remove {
			s.RemoveCode(id, scene)
		}
		return code, code != ""
	}
}

func (s *RedisCodeStore) RemoveCode(id string, scene string) {
	s.redisClient.Execute("DEL", s.RedisKey(CODE_REDIS_KEY_PREFIX, id, scene))
}

func (s *RedisCodeStore) RedisKey(name string, key string, scene string) string {
	return s.keyPrefix + name + key + ":" + scene
}
