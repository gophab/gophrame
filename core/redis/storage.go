package redis

import (
	"strings"
	"sync"
	"time"

	"github.com/wjshen/gophrame/core/json"

	"github.com/gomodule/redigo/redis"
	"github.com/timandy/routine"
)

type RedisLock struct {
	*RedisClient
	lockKey   string
	expireIn  time.Duration
	timeout   time.Duration
	locked    bool
	reentrant bool // 同线程可重入机制
	counter   routine.ThreadLocal
}

func NewRedisLock(client *RedisClient, key string) *RedisLock {
	return &RedisLock{
		RedisClient: client,
		lockKey:     key,
		expireIn:    60 * time.Second,
		timeout:     10 * time.Second,
		locked:      false,
		reentrant:   false,
	}
}

func (l *RedisLock) Locked() bool {
	return l.locked
}

func (l *RedisLock) Lock(timeout ...int) (bool, error) {
	if len(timeout) > 0 {
		l.timeout = time.Duration(timeout[0]) * time.Second
	}

	for t := l.timeout.Microseconds(); t >= 0; t -= 100 {
		if l.TryLock() {
			return true, nil
		}
	}
	return false, nil
}

func (l *RedisLock) Unlock() {
	if l.reentrant && l.counter.Get() != nil {
		c := *l.counter.Get().(*int) - 1
		l.counter.Set(&c)
		if c > 0 {
			return
		}
	}

	if l.locked {
		l.client.Do("DEL", l.lockKey)
		l.locked = false
	}
	l.counter.Remove()
}

func (l *RedisLock) TryLock() bool {
	if l.reentrant && l.counter.Get() != nil {
		c := *l.counter.Get().(*int) + 1
		l.counter.Set(&c)
	}

	expires := time.Now().Add(l.expireIn).Add(time.Second).UnixMilli()
	if b, err := l.client.Do(`SETNX`, l.lockKey, expires); err == nil && b.(bool) {
		l.locked = true
		c := 1
		l.counter.Set(&c)
		return true
	}

	if res, err := l.Int64(l.client.Do("GET", l.lockKey)); err == nil && res <= time.Now().UnixMilli() {
		// 已过期
		if res, err := l.Int64(l.client.Do("GETSET", l.lockKey, expires)); err == nil && res == expires {
			l.locked = true
			c := 1
			l.counter.Set(&c)
			return true
		}
	}

	return false
}

type RedisStorage struct {
	*RedisClient
	locks sync.Map
	mutex sync.Mutex
}

func NewRedisStorage(database int) *RedisStorage {
	return &RedisStorage{
		RedisClient: GetOneRedisClientIndex(database),
	}
}

func (s *RedisStorage) getLock(key string) *RedisLock {
	if result, b := s.locks.Load(key); b {
		return result.(*RedisLock)
	} else {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if result, b := s.locks.Load(key); b {
			return result.(*RedisLock)
		} else {
			result := NewRedisLock(s.RedisClient, key)
			s.locks.Store(key, result)
			return result
		}
	}
}

func (s *RedisStorage) Lock(key string) {
	s.getLock(key).Lock()
}

func (s *RedisStorage) TryLock(key string) bool {
	lock := s.getLock(key)
	lock.reentrant = false
	return lock.TryLock()
}

func (s *RedisStorage) ReenLock(key string) {
	lock := s.getLock(key)
	lock.reentrant = true
	lock.Lock()
}

func (s *RedisStorage) TryReenLock(key string) bool {
	lock := s.getLock(key)
	lock.reentrant = true
	return lock.TryLock()
}

func (s *RedisStorage) Unlock(key string) {
	lock := s.getLock(key)
	if lock != nil && lock.Locked() {
		lock.Unlock()
	}
}

func (s *RedisStorage) Restore(key string) *RedisResult {
	res, err := s.client.Do("GET", key)
	return s.Result(key, res, err)
}

func (s *RedisStorage) RestoreKeys(keys []string) *RedisResults {
	var args []interface{}
	for _, v := range keys {
		args = append(args, v)
	}
	res, err := s.client.Do("MGET", args...)
	return s.Results(keys, res, err)
}

func (s *RedisStorage) Store(key string, value interface{}) error {
	_, err := s.client.Do("SET", key, value)
	return err
}

func (s *RedisStorage) StoreDuration(key string, value interface{}, duration time.Duration) error {
	_, err := s.client.Do("PSETEX", key, duration.Milliseconds(), value)
	if err != nil {
		_, err = s.client.Do("SETEX", key, duration.Seconds(), value)
	}
	return err
}

func (s *RedisStorage) StoreExpire(key string, value interface{}, expireTime time.Time) error {
	_, err := s.client.Do("PSETEX", key, time.Until(expireTime).Microseconds(), value)
	if err != nil {
		s.client.Do("SETEX", key, time.Until(expireTime).Seconds(), value)
	}
	return err
}

func (s *RedisStorage) StoreIfAbsent(key string, value interface{}) (string, error) {
	res, err := s.client.Do("SETNX", key, value)
	return s.String(res, err)
}

func (s *RedisStorage) StoreDurationIfAbsent(key string, value interface{}, duration time.Duration) (string, error) {
	defer func() {
		s.client.Do("PEXPIRE", key, duration.Milliseconds())
	}()
	res, err := s.client.Do("SETNX", key, value)
	return s.String(res, err)
}

func (s *RedisStorage) StoreAsString(key string, value interface{}) error {
	_, err := s.client.Do("SET", key, json.String(value))
	return err
}

func (s *RedisStorage) StoreAsStringDuration(key string, value interface{}, duration time.Duration) error {
	_, err := s.client.Do("PSETEX", key, duration.Milliseconds(), json.String(value))
	if err != nil {
		_, err = s.client.Do("SETEX", key, duration.Seconds(), json.String(value))
	}
	return err
}

func (s *RedisStorage) StoreAsStringExpire(key string, value interface{}, expireTime time.Time) error {
	_, err := s.client.Do("PSETEX", key, time.Until(expireTime).Microseconds(), json.String(value))
	if err != nil {
		s.client.Do("SETEX", key, time.Until(expireTime).Seconds(), json.String(value))
	}
	return err
}

func (s *RedisStorage) StoreAsStringIfAbsent(key string, value interface{}) (string, error) {
	res, err := s.client.Do("SETNX", key, json.String(value))
	return s.String(res, err)
}

func (s *RedisStorage) StoreAsStringDurationIfAbsent(key string, value interface{}, duration time.Duration) (string, error) {
	defer func() {
		s.client.Do("PEXPIRE", key, duration.Milliseconds())
	}()
	res, err := s.client.Do("SETNX", key, json.String(value))
	return s.String(res, err)
}

func (s *RedisStorage) StoreSetValue(key string, value interface{}) (int, error) {
	res, err := s.client.Do("SADD", key, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreSetValueAsString(key string, value interface{}) (int, error) {
	res, err := s.client.Do("SADD", key, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) StoreSetValues(key string, values []interface{}) (int, error) {
	p := []interface{}{key, values[0:]}
	resp, err := s.client.Do("SADD", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) StoreSetValuesAsString(key string, values []interface{}) (int, error) {
	p := []interface{}{key}
	for _, v := range values {
		p = append(p, json.String(v))
	}
	resp, err := s.client.Do("SADD", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) RestoreSet(key string) *RedisResult {
	res, err := s.client.Do("SMEMBERS", key)
	return s.Result(key, res, err)
}

func (s *RedisStorage) RestoreSets(keys []string) *RedisResult {
	p := []interface{}{keys[0:]}
	res, err := s.client.Do("SUNION", p)
	return s.Result(strings.Join(keys, ","), res, err)
}

func (s *RedisStorage) StoreZSetValue(key string, score int, value interface{}) (int, error) {
	res, err := s.client.Do("ZADD", key, score, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreZSetValueAsString(key string, score int, value interface{}) (int, error) {
	res, err := s.client.Do("ZADD", key, score, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) RestoreZSet(key string) *RedisResult {
	res, err := s.client.Do("ZRANGE", key, 0, 0)
	return s.Result(key, res, err)
}

func (s *RedisStorage) StoreListValue(key string, value interface{}) (int, error) {
	res, err := s.client.Do("LPUSH", key, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValueTail(key string, value interface{}) (int, error) {
	res, err := s.client.Do("LPUSH", key, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValueHead(key string, value interface{}) (int, error) {
	res, err := s.client.Do("RPUSH", key, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValueAsString(key string, value interface{}) (int, error) {
	res, err := s.client.Do("LPUSH", key, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValueTailAsString(key string, value interface{}) (int, error) {
	res, err := s.client.Do("LPUSH", key, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValueHeadAsString(key string, value interface{}) (int, error) {
	res, err := s.client.Do("RPUSH", key, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) StoreListValues(key string, values []interface{}) (int, error) {
	p := []interface{}{key, values[0:]}
	resp, err := s.client.Do("LPUSH", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) StoreListValuesAsString(key string, values []interface{}) (int, error) {
	p := []interface{}{key}
	for _, v := range values {
		p = append(p, json.String(v))
	}
	resp, err := s.client.Do("LPUSH", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) RestoreList(key string) *RedisResult {
	res, err := s.client.Do("LRANGE", key, 0, 0)
	return s.Result(key, res, err)
}

func (s *RedisStorage) RestoreListHead(key string) *RedisResult {
	res, err := s.client.Do("RPOP", key)
	return s.Result(key, res, err)
}

func (s *RedisStorage) RestoreListTail(key string) *RedisResult {
	res, err := s.client.Do("LPOP", key)
	return s.Result(key, res, err)
}

func (s *RedisStorage) StoreHashValue(key string, hashKey string, value interface{}) (int, error) {
	res, err := s.client.Do("HSET", key, hashKey, value)
	return s.Int(res, err)
}

func (s *RedisStorage) StoreHashValueAsString(key string, value interface{}) (int, error) {
	res, err := s.client.Do("HSET", key, json.String(value))
	return s.Int(res, err)
}

func (s *RedisStorage) StoreHashValues(key string, values map[string]interface{}) (int, error) {
	p := []interface{}{key}
	for k, v := range values {
		p = append(p, k, v)
	}
	resp, err := s.client.Do("HMSET", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) StoreHashValuesAsString(key string, values map[string]interface{}) (int, error) {
	p := []interface{}{key}
	for k, v := range values {
		p = append(p, k, json.String(v))
	}
	resp, err := s.client.Do("HMSET", p)
	return s.Int(resp, err)
}

func (s *RedisStorage) RestoreHash(key string) *RedisResult {
	keys, err := redis.Strings(s.client.Do("HKEYS", key))
	if err != nil {
		return s.Result(key, nil, err)
	}
	values, err := redis.Values(s.client.Do("HVALS", key))
	if err != nil {
		return s.Result(key, nil, err)
	}

	res := make(map[string]interface{})
	for i := 0; i < len(keys); i++ {
		res[keys[i]] = values[i]
	}
	return s.Result(key, res, err)
}

func (s *RedisStorage) RestoreHashValue(key string, hashKey string) *RedisResult {
	res, err := s.client.Do("HGET", key, hashKey)
	return s.Result(key+":"+hashKey, res, err)
}
