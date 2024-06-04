package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gophab/gophrame/core/eventbus"
	"github.com/gophab/gophrame/core/global"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/redis/config"

	"github.com/gomodule/redigo/redis"
)

const (
	//redis部分
	ErrorsRedisInitConnFail string = "初始化redis连接池失败"
	ErrorsRedisAuthFail     string = "Redis Auth 鉴权失败，密码错误"
	ErrorsRedisGetConnFail  string = "Redis 从连接池获取一个连接失败，超过最大重试次数"
)

var redisPools map[string]*redis.Pool = make(map[string]*redis.Pool)

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func initRedisClientPool(databaseIndex int) *redis.Pool {
	if result := redisPools[strconv.Itoa(databaseIndex)]; result != nil {
		return result
	}

	result := &redis.Pool{
		MaxIdle:     config.Setting.MaxIdle,    //最大空闲数
		MaxActive:   config.Setting.MaxActive,  //最大活跃数
		IdleTimeout: config.Setting.IdleTimout, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			//此处对应redis ip及端口号
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Setting.Host, config.Setting.Port))
			if err != nil {
				logger.Error(ErrorsRedisInitConnFail, err.Error())
				return nil, err
			}
			auth := config.Setting.Auth //通过配置项设置redis密码
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					_ = conn.Close()
					logger.Error(ErrorsRedisAuthFail, err.Error())
				}
			}
			_, _ = conn.Do("select", databaseIndex)
			return conn, err
		},
	}

	// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
	eventbus.RegisterEventListener(global.EventDestroyPrefix+"Redis"+strconv.Itoa(databaseIndex), func(event string, args ...interface{}) {
		_ = result.Close()
	})

	redisPools[strconv.Itoa(databaseIndex)] = result

	return result
}

// 从连接池获取一个redis连接
func GetOneRedisClient() *RedisClient {
	return GetOneRedisClientIndex(config.Setting.Database)
}

func GetOneRedisClientIndex(databaseIndex int) *RedisClient {
	maxRetryTimes := config.Setting.ConnectionFailRetryTimes
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = initRedisClientPool(databaseIndex).Get()
		// 首先通过执行一个获取时间的命令检测连接是否有效，如果已有的连接无法执行命令，则重新尝试连接到redis服务器获取新的连接池地址
		// 连接不可用可能会发生的场景主要有：服务端redis重启、客户端网络在有线和无线之间切换等
		if _, replyErr := oneConn.Do("time"); replyErr != nil {
			//fmt.Printf("连接已经失效(出错)：%+v\n", replyErr.Error())
			// 如果已有的redis连接池获取连接出错(官方库的说法是连接不可用)，那么继续使用从新初始化连接池
			redisPools[strconv.Itoa(databaseIndex)] = nil
			oneConn = initRedisClientPool(databaseIndex).Get()
		}
		if oneConn.Err() != nil {
			if i == maxRetryTimes {
				logger.Error(ErrorsRedisGetConnFail, oneConn.Err())
				return nil
			}
			//如果出现网络短暂的抖动，短暂休眠后，支持自动重连
			time.Sleep(config.Setting.ReConnectInterval)
		} else {
			break
		}
	}
	return &RedisClient{oneConn}
}

// 定义一个redis客户端结构体
type RedisClient struct {
	client redis.Conn
}

// 为redis-go 客户端封装统一操作函数入口
func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

// 释放连接到连接池
func (r *RedisClient) ReleaseOneRedisClient() {
	_ = r.client.Close()
}

// bool 类型转换
func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

// string 类型转换
func (r *RedisClient) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// strings 类型转换
func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

// Float64 类型转换
func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

// int 类型转换
func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// int64 类型转换
func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// uint64 类型转换
func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

// Bytes 类型转换
func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

type RedisResult struct {
	key   string
	reply interface{}
	err   error
}

// 封装几个数据类型转换的函数
func (r *RedisClient) Result(key string, reply interface{}, err error) *RedisResult {
	return &RedisResult{
		key,
		reply,
		err,
	}
}

// 封装几个数据类型转换的函数
func (r *RedisClient) Results(keys []string, reply interface{}, err error) *RedisResults {
	results := make([]*RedisResult, 0)
	if values, err := redis.Values(reply, err); err == nil {
		for i := 0; i < len(values); i++ {
			results = append(results, &RedisResult{
				key:   keys[i],
				reply: values[i],
				err:   err,
			})
		}
	}
	return &RedisResults{
		keys,
		reply,
		err,
		results,
	}
}

// bool 类型转换
func (r *RedisResult) Bool() (bool, error) {
	return redis.Bool(r.reply, r.err)
}

// string 类型转换
func (r *RedisResult) String() (string, error) {
	return redis.String(r.reply, r.err)
}

// strings 类型转换
func (r *RedisResult) Strings() ([]string, error) {
	return redis.Strings(r.reply, r.err)
}

// Float64 类型转换
func (r *RedisResult) Float64() (float64, error) {
	return redis.Float64(r.reply, r.err)
}

// int 类型转换
func (r *RedisResult) Int() (int, error) {
	return redis.Int(r.reply, r.err)
}

// int64 类型转换
func (r *RedisResult) Int64() (int64, error) {
	return redis.Int64(r.reply, r.err)
}

// uint64 类型转换
func (r *RedisResult) Uint64() (uint64, error) {
	return redis.Uint64(r.reply, r.err)
}

// Bytes 类型转换
func (r *RedisResult) Bytes() ([]byte, error) {
	return redis.Bytes(r.reply, r.err)
}

func (r *RedisResult) As(data interface{}) error {
	if rep, err := r.String(); err == nil {
		return json.Unmarshal([]byte(rep), data)
	} else {
		return err
	}
}

type RedisResults struct {
	keys    []string
	reply   interface{}
	err     error
	results []*RedisResult
}

func (rs *RedisResults) List() []*RedisResult {
	return rs.results
}

func (rs *RedisResults) Map() map[string]*RedisResult {
	result := make(map[string]*RedisResult)
	for _, rr := range rs.results {
		result[rr.key] = rr
	}
	return result
}

func (rs *RedisResults) AsStringList() []string {
	result := make([]string, 0)
	for _, rr := range rs.results {
		if s, err := rr.String(); err == nil {
			result = append(result, s)
		}
	}
	return result
}

func (rs *RedisResults) AsStringMap() map[string]string {
	result := make(map[string]string)
	for _, rr := range rs.results {
		if s, err := rr.String(); err == nil {
			result[rr.key] = s
		}
	}
	return result
}
