package token

import (
	"context"
	"errors"
	"time"

	"github.com/gophab/gophrame/core/database"
	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/json"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/redis"
	"github.com/gophab/gophrame/core/security/token/config"
	"github.com/gophab/gophrame/core/util"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	theTokenStore oauth2.TokenStore
)

func TokenStore() oauth2.TokenStore {
	if theTokenStore == nil {
		theTokenStore = InitTokenStore()
	}
	return theTokenStore
}

func InitTokenStore() oauth2.TokenStore {
	if theTokenStore == nil {
		var store oauth2.TokenStore
		var err error

		switch config.Setting.Store.Mode {
		case "database":
			store, err = NewDatabaseTokenStore()
		case "redis":
			store, err = NewRedisTokenStore()
		case "file":
			store, err = NewFileTokenStore(config.Setting.Store.File.FileName)
		default:
			store, err = NewMemeoryTokenStore()
		}

		if err == nil && store != nil {
			inject.InjectValue("tokenStore", store)
		}

		theTokenStore = store
	}

	return theTokenStore
}

func NewMemeoryTokenStore() (oauth2.TokenStore, error) {
	logger.Debug("Using memory token store")
	return store.NewMemoryTokenStore()
}

func NewFileTokenStore(filename string) (oauth2.TokenStore, error) {
	logger.Debug("Using file token store")
	return store.NewFileTokenStore(filename)
}

type ITokenStore interface {
	oauth2.TokenStore
	GetToken(context.Context, string) (oauth2.TokenInfo, error)
}

type DatabaseTokenStore struct {
}

func NewDatabaseTokenStore() (oauth2.TokenStore, error) {
	logger.Debug("Using database token store")
	result := &DatabaseTokenStore{}
	go func() {
		for {
			result.clearExpiredTokens()

			// 延时30秒
			time.Sleep(time.Second * 30)
		}
	}()
	return result, nil
}

func (s *DatabaseTokenStore) clearExpiredTokens() {
	now := time.Now()
	database.DB().Exec(`DELETE FROM oauth_access_token WHERE expiration < ?`, now)
	database.DB().Exec(`DELETE FROM oauth_refresh_token WHERE expiration < ?`, now)
	database.DB().Exec(`DELETE FROM oauth_code WHERE expiration < ?`, now)
}

func (s *DatabaseTokenStore) insertAccessToken(info oauth2.TokenInfo) (int64, error) {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}
	result := database.DB().Exec(
		`INSERT INTO oauth_access_token (access_token, token, authentication_id, authentication, client_id, user_name, refresh_token, expiration) 
		 SELECT ?, ?, ?, ?, ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM oauth_access_token WHERE authentication_id=?)`,
		util.MD5(info.GetAccess()),
		json.String(info),
		authentication.GetId(),
		json.String(authentication),
		info.GetClientID(),
		info.GetUserID(),
		util.MD5(info.GetRefresh()),
		info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()),
		authentication.GetId())
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *DatabaseTokenStore) insertRefreshToken(info oauth2.TokenInfo) (int64, error) {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}
	result := database.DB().Exec(
		`INSERT INTO oauth_refresh_token (refresh_token, token, authentication, expiration) 
		 SELECT ?, ?, ?, ? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM oauth_refresh_token WHERE refresh_token=?)`,
		util.MD5(info.GetRefresh()),
		json.String(info),
		json.String(authentication),
		info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()),
		util.MD5(info.GetRefresh()))
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (s *DatabaseTokenStore) insertCode(info oauth2.TokenInfo) (int64, error) {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}
	result := database.DB().Exec(
		`INSERT INTO oauth_code (code, authentication, expiration) 
		 SELECT ?, ?, ? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM oauth_code WHERE code=?)`,
		util.MD5(info.GetCode()),
		json.String(authentication),
		info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()),
		util.MD5(info.GetCode()))
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// create and store the new token information
func (s *DatabaseTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}

	if info.GetAccess() != "" {
		if rows, err := s.insertAccessToken(info); err != nil {
			return err
		} else if rows <= 0 {
			// 已存在授权
			// 1. 删除已有的refresh_token
			// database.DB().Exec(`DELETE FROM oauth_refresh_token WHERE refresh_token IN (SELECT refresh_token FROM oauth_access_token WHERE authentication_id=?)`, authentication.GetId())

			// 2. 获取现有
			if exist, _ := s.GetToken(ctx, authentication.GetId()); exist != nil {
				// 现有记录
				if info.GetRefresh() == exist.GetRefresh() {
					info.SetRefreshCreateAt(exist.GetRefreshCreateAt())
					info.SetRefreshExpiresIn(exist.GetRefreshExpiresIn())
					// info.SetRefreshExpiresIn(time.Until(exist.GetRefreshCreateAt().Add(exist.GetRefreshExpiresIn())))

					// 未更新RefreshToken
					if info.GetAccess() == exist.GetAccess() {
						info.SetAccessCreateAt(exist.GetAccessCreateAt())
						//info.SetAccessExpiresIn(exist.GetAccessExpiresIn())
						info.SetAccessExpiresIn(time.Until(exist.GetAccessCreateAt().Add(exist.GetAccessExpiresIn())))
						// reuse 已有的AccessToken，不更新token
					} else {
						exist.SetAccess(info.GetAccess())
						exist.SetAccessCreateAt(info.GetAccessCreateAt())
						exist.SetAccessExpiresIn(info.GetAccessExpiresIn())

						// reuse已有的RefreshToken，只更新access_token，
						if err := database.DB().Exec(
							`UPDATE oauth_access_token SET access_token=?, token=?, expiration=? WHERE authentication_id=?`,
							util.MD5(exist.GetAccess()),
							json.String(exist),
							info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()),
							authentication.GetId()).Error; err != nil {
							return err
						}
					}
				} else {
					// 完全更新
					if err := database.DB().Exec(
						`UPDATE oauth_access_token SET access_token=?, token=?, client_id=?, user_name=?, refresh_token=?, expiration=? WHERE authentication_id=?`,
						util.MD5(info.GetAccess()),
						json.String(info),
						info.GetClientID(),
						info.GetUserID(),
						util.MD5(info.GetRefresh()),
						info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()),
						authentication.GetId()).Error; err != nil {
						return err
					}
				}
			} else {
				s.insertAccessToken(info)
			}
		}
	}

	if info.GetRefresh() != "" {
		if rows, err := s.insertRefreshToken(info); err != nil {
			return err
		} else if rows <= 0 {
			if exist, _ := s.GetByRefresh(ctx, authentication.GetId()); exist != nil {
				if err := database.DB().Exec(
					`UPDATE oauth_refresh_token SET token=?, authentication=?, expirtaion=? WHERE refresh_token=?`,
					json.String(info),
					json.String(authentication),
					info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()),
					util.MD5(info.GetRefresh())).Error; err != nil {
					return err
				}
			} else {
				s.insertRefreshToken(info)
			}
		}
	}

	if info.GetCode() != "" {
		if rows, err := s.insertCode(info); err != nil {
			return err
		} else if rows <= 0 {
			if exist, _ := s.GetByCode(ctx, info.GetCode()); exist != nil {
				if err := database.DB().Exec(
					`UPDATE oauth_code SET authentication=?, expiration=? WHERE code=?`,
					json.String(authentication),
					info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()),
					util.MD5(info.GetCode())).Error; err != nil {
					return err
				}
			} else {
				s.insertCode(info)
			}
		}
	}

	return nil
}

// delete the authorization code
func (s *DatabaseTokenStore) RemoveByCode(ctx context.Context, code string) error {
	tokenInfo, err := s.GetByCode(ctx, code)
	if err != nil {
		return err
	}

	if tokenInfo != nil {
		s.RemoveByAccess(ctx, tokenInfo.GetAccess())
		s.RemoveByRefresh(ctx, tokenInfo.GetRefresh())

		return database.DB().Exec("DELETE FROM oauth_code WHERE code = ?", util.MD5(code)).Error
	}

	return errors.New("not found")
}

// use the access token to delete the token information
func (s *DatabaseTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return database.DB().Exec("DELETE FROM oauth_access_token WHERE access_token = ?", util.MD5(access)).Error
}

// use the refresh token to delete the token information
func (s *DatabaseTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	if tokenInfo, err := s.GetByRefresh(ctx, refresh); err == nil {
		database.DB().Exec("DELETE FROM oauth_access_token WHERE access_token = ?", util.MD5(tokenInfo.GetAccess()))
		return database.DB().Exec("DELETE FROM oauth_refresh_token WHERE refresh_token = ?", util.MD5(refresh)).Error
	}
	return nil
}

// use the authorization code for token information data
func (s *DatabaseTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	var authenticationString string
	result := database.DB().Raw("SELECT authentication FROM oauth_code WHERE code = ? LIMIT 1", util.MD5(code)).First(&authenticationString)
	if result.Error != nil || result.RowsAffected <= 0 {
		return nil, result.Error
	}

	authenticaiton, err := ParseAuthentication(authenticationString)
	if err != nil {
		return nil, err
	}

	var tokenString string
	result = database.DB().Raw("SELECT token FROM oauth_access_token WHERE authentication_id = ? LIMIT 1", authenticaiton.GetId()).First(&tokenString)
	if result.Error == nil && result.RowsAffected > 0 {
		return ParseToken(tokenString)
	}
	return nil, result.Error
}

// use the access token for token information data
func (s *DatabaseTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	var tokenString string
	result := database.DB().Raw("SELECT token FROM oauth_access_token WHERE access_token = ? LIMIT 1", util.MD5(access)).First(&tokenString)
	if result.Error == nil && result.RowsAffected > 0 {
		return ParseToken(tokenString)
	}
	return nil, result.Error
}

// use the refresh token for token information data
func (s *DatabaseTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	var tokenString string
	result := database.DB().Raw("SELECT token FROM oauth_refresh_token WHERE refresh_token = ? LIMIT 1", util.MD5(refresh)).First(&tokenString)
	if result.Error == nil && result.RowsAffected > 0 {
		return ParseToken(tokenString)
	}
	return nil, result.Error
}

func (s *DatabaseTokenStore) GetToken(ctx context.Context, key string) (oauth2.TokenInfo, error) {
	var tokenString string
	result := database.DB().Raw("SELECT token FROM oauth_access_token WHERE authentication_id = ? LIMIT 1", key).First(&tokenString)
	if result.Error == nil && result.RowsAffected > 0 {
		return ParseToken(tokenString)
	}
	return nil, result.Error
}

func ParseToken(tokenString string) (oauth2.TokenInfo, error) {
	var token OAuth2Token
	json.Json(tokenString, &token)
	return &token, nil
}

func ParseAuthentication(authenticationString string) (*Authentication, error) {
	var authentication Authentication
	json.Json(authenticationString, &authentication)
	return &authentication, nil
}

/**
 * Redis TokenStore: store access_token/refresh_token into redis
 */

const (
	ACCESS              = "access:"
	AUTH_TO_ACCESS      = "auth_to_access:"
	AUTH                = "auth:"
	REFRESH_AUTH        = "refresh_auth:"
	ACCESS_TO_REFRESH   = "access_to_refresh:"
	REFRESH             = "refresh:"
	REFRESH_TO_ACCESS   = "refresh_to_access:"
	CLIENT_ID_TO_ACCESS = "client_id_to_access:"
	UNAME_TO_ACCESS     = "uname_to_access:"
)

type RedisTokenStore struct {
	redisClient *redis.RedisClient
	keyPrefix   string
}

func NewRedisTokenStore() (oauth2.TokenStore, error) {
	logger.Debug("Using redis token store")
	var tokenStore oauth2.TokenStore
	if config.Setting.Store.Redis != nil {
		tokenStore = &RedisTokenStore{
			redisClient: redis.GetOneRedisClientIndex(config.Setting.Store.Redis.Database),
			keyPrefix:   config.Setting.Store.Redis.KeyPrefix,
		}
	}
	return tokenStore, nil
}

func (s *RedisTokenStore) cacheAuthentication(ctx context.Context, info oauth2.TokenInfo) error {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(AUTH, info.GetAccess()), info.GetAccessExpiresIn().Seconds(), json.String(authentication)); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(REFRESH_AUTH, info.GetRefresh()), info.GetRefreshExpiresIn().Seconds(), json.String(authentication)); err != nil {
		return err
	}

	return nil
}

func (s *RedisTokenStore) cacheAccess(ctx context.Context, info oauth2.TokenInfo) error {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(AUTH_TO_ACCESS, authentication.GetId()), info.GetAccessExpiresIn().Seconds(), util.MD5(info.GetAccess())); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("ZADD", s.RedisKey(CLIENT_ID_TO_ACCESS, info.GetClientID()), time.Now().Nanosecond(), util.MD5(info.GetAccess())); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("ZADD", s.RedisKey(UNAME_TO_ACCESS, info.GetUserID()), time.Now().Nanosecond(), util.MD5(info.GetAccess())); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(ACCESS_TO_REFRESH, util.MD5(info.GetAccess())), info.GetAccessExpiresIn().Seconds(), util.MD5(info.GetRefresh())); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(REFRESH_TO_ACCESS, util.MD5(info.GetRefresh())), info.GetAccessExpiresIn().Seconds(), util.MD5(info.GetAccess())); err != nil {
		return err
	}

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(ACCESS, util.MD5(info.GetAccess())), info.GetAccessExpiresIn().Seconds(), json.String(info)); err != nil {
		return err
	}

	return nil
}

func (s *RedisTokenStore) cacheRefresh(ctx context.Context, info oauth2.TokenInfo) error {

	if _, err := s.redisClient.Execute("SETEX", s.RedisKey(REFRESH, util.MD5(info.GetRefresh())), info.GetRefreshExpiresIn().Seconds(), json.String(info)); err != nil {
		return err
	}

	return nil
}

// create and store the new token information
func (s *RedisTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	var authentication Authentication = Authentication{
		UserId:        info.GetUserID(),
		Authenticated: true,
		Request: TokenRequest{
			ClientId: info.GetClientID(),
			Scope:    info.GetScope(),
		},
	}

	exist, _ := s.GetToken(ctx, authentication.GetId())
	if exist != nil && exist.GetRefresh() == info.GetRefresh() {
		info.SetRefreshCreateAt(exist.GetRefreshCreateAt())
		info.SetRefreshExpiresIn(exist.GetRefreshExpiresIn())
		// info.SetRefreshExpiresIn(time.Until(exist.GetRefreshCreateAt().Add(exist.GetRefreshExpiresIn())))

		if exist.GetAccess() == info.GetAccess() {
			info.SetAccessCreateAt(exist.GetAccessCreateAt())
			//info.SetAccessExpiresIn(exist.GetAccessExpiresIn())
			info.SetAccessExpiresIn(time.Until(exist.GetAccessCreateAt().Add(exist.GetAccessExpiresIn())))
			// reuse 已有的AccessToken，不更新token
		} else {
			exist.SetAccess(info.GetAccess())
			exist.SetAccessCreateAt(info.GetAccessCreateAt())
			exist.SetAccessExpiresIn(info.GetAccessExpiresIn())

			// reuse已有的RefreshToken，只更新access_token，
			if err := s.cacheAccess(ctx, exist); err != nil {
				return err
			}
		}
	} else {
		s.cacheAuthentication(ctx, info)
		s.cacheAccess(ctx, info)
		s.cacheRefresh(ctx, info)
	}

	return nil
}

// delete the authorization code
func (s *RedisTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return nil
}

// use the access token to delete the token information
func (s *RedisTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	info, _ := s.GetByAccess(ctx, access)
	if info != nil {
		var authentication Authentication = Authentication{
			UserId:        info.GetUserID(),
			Authenticated: true,
			Request: TokenRequest{
				ClientId: info.GetClientID(),
				Scope:    info.GetScope(),
			},
		}

		if _, err := s.redisClient.Execute("DEL", s.RedisKey(AUTH_TO_ACCESS, authentication.GetId())); err != nil {
			return err
		}

		if _, err := s.redisClient.Execute("ZREM", s.RedisKey(CLIENT_ID_TO_ACCESS, info.GetClientID()), util.MD5(access)); err != nil {
			return err
		}

		if _, err := s.redisClient.Execute("ZREM", s.RedisKey(UNAME_TO_ACCESS, info.GetUserID()), util.MD5(access)); err != nil {
			return err
		}

		if _, err := s.redisClient.Execute("DEL", s.RedisKey(ACCESS_TO_REFRESH, util.MD5(access))); err != nil {
			return err
		}

		if _, err := s.redisClient.Execute("DEL", s.RedisKey(ACCESS, util.MD5(access))); err != nil {
			return err
		}
	}

	return nil
}

// use the refresh token to delete the token information
func (s *RedisTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	tokenInfo, _ := s.GetByRefresh(ctx, refresh)
	if tokenInfo != nil {
		s.RemoveByAccess(ctx, tokenInfo.GetAccess())

		// Remove REFRESH
		if _, err := s.redisClient.Execute("DEL", s.RedisKey(REFRESH_TO_ACCESS, util.MD5(refresh))); err != nil {
			return err
		}
		if _, err := s.redisClient.Execute("DEL", s.RedisKey(REFRESH, util.MD5(refresh))); err != nil {
			return err
		}
		if _, err := s.redisClient.Execute("DEL", s.RedisKey(REFRESH_AUTH, util.MD5(refresh))); err != nil {
			return err
		}
	}

	return nil
}

// use the authorization code for token information data
func (s *RedisTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	// AUTH_TO_ACCESS => accessToken => ACCESS => TokenString
	if md5AccessToken, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(AUTH_TO_ACCESS, code))); err != nil {
		return nil, err
	} else if tokenString, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(ACCESS, md5AccessToken))); err != nil {
		return nil, err
	} else {
		return ParseToken(tokenString)
	}
}

// use the access token for token information data
func (s *RedisTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	// ACCESS => TokenString
	if tokenString, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(ACCESS, util.MD5(access)))); err != nil {
		return nil, err
	} else {
		return ParseToken(tokenString)
	}
}

// use the refresh token for token information data
func (s *RedisTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	// REFRESH => TokenString
	if tokenString, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(REFRESH, util.MD5(refresh)))); err != nil {
		return nil, err
	} else {
		return ParseToken(tokenString)
	}
}

func (s *RedisTokenStore) GetToken(ctx context.Context, authenticationId string) (oauth2.TokenInfo, error) {
	// AUTH_TO_ACCESS => accessToken => ACCESS => TokenString
	if accessToken, err := s.redisClient.String(s.redisClient.Execute("GET", s.RedisKey(AUTH_TO_ACCESS, authenticationId))); err != nil {
		return nil, err
	} else {
		return s.GetByAccess(ctx, accessToken)
	}
}

func (s *RedisTokenStore) RedisKey(name string, key string) string {
	return s.keyPrefix + name + key
}
