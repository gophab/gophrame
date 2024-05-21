package token

import (
	"strconv"
	"strings"
	"time"

	"github.com/wjshen/gophrame/core/database"
	"github.com/wjshen/gophrame/core/global"
	"github.com/wjshen/gophrame/core/logger"
	"github.com/wjshen/gophrame/core/redis"
	"github.com/wjshen/gophrame/core/security/token/config"
	"github.com/wjshen/gophrame/core/util"
)

// 本文件专门处理 token 缓存到 redis 的相关逻辑
func ValidTokenCacheToRedis(userId int64) {
	tokenCacheRedisFact := CreateRedisUserTokenStorage(userId)
	if tokenCacheRedisFact == nil {
		logger.Error("redis连接失败，请检查配置")
		return
	}
	defer tokenCacheRedisFact.ReleaseRedisConn()

	sql := "SELECT   token,expires_at  FROM  `tb_auth_access_tokens`  WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
	maxOnlineUsers := config.Setting.OnlineUsers
	rows, err := database.DB().Raw(sql, userId, maxOnlineUsers).Rows()
	defer func() {
		//  凡是获取原生结果集的查询，记得释放记录集
		_ = rows.Close()
	}()

	var tempToken, expires string
	if err == nil && rows != nil {
		for i := 1; rows.Next(); i++ {
			err = rows.Scan(&tempToken, &expires)
			if err == nil {
				if ts, err := time.ParseInLocation(global.DateFormat, expires, time.Local); err == nil {
					tokenCacheRedisFact.SetToken(ts.Unix(), tempToken)
					// 因为每个用户的token是按照过期时间倒叙排列的，第一个是有效期最长的，将该用户的总键设置一个最大过期时间，到期则自动清理，避免不必要的数据残留
					if i == 1 {
						tokenCacheRedisFact.SetUserTokensExpire(ts.Unix())
					}
				} else {
					logger.Error("expires_at 转换位时间戳出错", err.Error())
				}
			}
		}
	}
	// 缓存结束之后删除超过系统设置最大在线数量的token
	tokenCacheRedisFact.DelOverMaxOnlineCache()
}

// DelTokenCacheFromRedis 用户密码修改后，删除redis所有的token
func DelTokenCacheFromRedis(userId int64) {
	tokenCacheRedisFact := CreateRedisUserTokenStorage(userId)
	if tokenCacheRedisFact == nil {
		logger.Error("redis连接失败，请检查配置")
		return
	}
	tokenCacheRedisFact.ClearUserTokens()
	tokenCacheRedisFact.ReleaseRedisConn()
}

func CreateRedisUserTokenStorage(userId int64) *RedisUserTokenCache {
	client := redis.GetOneRedisClient()
	if client == nil {
		return nil
	}
	return &RedisUserTokenCache{redisClient: client, userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10)}
}

// 保存用户的所有Token
type RedisUserTokenCache struct {
	redisClient  *redis.RedisClient
	userTokenKey string
}

// SetTokenCache 设置缓存
func (u *RedisUserTokenCache) SetToken(tokenExpire int64, token string) bool {
	// 存储用户token时转为MD5，下一步比较的时候可以更加快速地比较是否一致
	if _, err := u.redisClient.Int(u.redisClient.Execute("zAdd", u.userTokenKey, tokenExpire, util.MD5(token))); err == nil {
		return true
	}
	return false
}

// DelOverMaxOnlineCache 删除缓存,删除超过系统允许最大在线数量之外的用户
func (u *RedisUserTokenCache) DelOverMaxOnlineCache() bool {
	// 首先先删除过期的token
	_, _ = u.redisClient.Execute("zRemRangeByScore", u.userTokenKey, 0, time.Now().Unix()-1)

	onlineUsers := config.Setting.OnlineUsers
	alreadyCacheNum, err := u.redisClient.Int(u.redisClient.Execute("zCard", u.userTokenKey))
	if err == nil && alreadyCacheNum > onlineUsers {
		// 删除超过最大在线数量之外的token
		if _, err = u.redisClient.Int(u.redisClient.Execute("zRemRangeByRank", u.userTokenKey, 0, alreadyCacheNum-onlineUsers-1)); err == nil {
			return true
		} else {
			logger.Error("删除超过系统允许之外的token出错：", err.Error())
		}
	}
	return false
}

// TokenCacheIsExists 查询token是否在redis存在
func (u *RedisUserTokenCache) UserTokenExists(token string) (exists bool) {
	token = util.MD5(token)

	curTimestamp := time.Now().Unix()
	onlineUsers := config.Setting.OnlineUsers
	if strSlice, err := u.redisClient.Strings(u.redisClient.Execute("zRevRange", u.userTokenKey, 0, onlineUsers-1)); err == nil {
		for _, val := range strSlice {
			if score, err := u.redisClient.Int64(u.redisClient.Execute("zScore", u.userTokenKey, token)); err == nil {
				if score > curTimestamp {
					if strings.Compare(val, token) == 0 {
						exists = true
						break
					}
				}
			}
		}
	} else {
		logger.Error("获取用户在redis缓存的 token 值出错：", err.Error())
	}
	return
}

// SetUserTokenExpire 设置用户的 usertoken 键过期时间
// 参数： 时间戳
func (u *RedisUserTokenCache) SetUserTokensExpire(ts int64) bool {
	if _, err := u.redisClient.Execute("expireAt", u.userTokenKey, ts); err == nil {
		return true
	}
	return false
}

// ClearUserToken 清除某个用户的全部缓存，当用户更改密码或者用户被禁用则删除该用户的全部缓存
func (u *RedisUserTokenCache) ClearUserTokens() bool {
	if _, err := u.redisClient.Execute("del", u.userTokenKey); err == nil {
		return true
	}
	return false
}

// ReleaseRedisConn 释放redis
func (u *RedisUserTokenCache) ReleaseRedisConn() {
	u.redisClient.ReleaseOneRedisClient()
}
