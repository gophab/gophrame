package server

import (
	"time"
)

/**
 * OAuth2 Client: webapp/1234567890
 */
type OAuthClient struct {
	ClientId     string    `gorm:"primaryKey"`
	ClientSecret string    ``
	ResourceIds  string    ``
	Scope        string    ``
	CreatedBy    string    `json:"created_by"`
	CreatedTime  time.Time `gorm:"autoCreateTime;type:TIMESTAMP;default:CURRENT_TIMESTAMP;<-:create" json:"created_time"`
	ModifiedBy   string    `json:"modified_by"`
	ModifiedTime time.Time `gorm:"autoUpdateTime;type:TIMESTAMP;default:CURRENT_TIMESTAMP on update current_timestamp" json:"modified_time"`
	DelFlag      bool      `gorm:"default:false" json:"del_flag"`
}

func (c *OAuthClient) TableName() string {
	return "oauth_client_details"
}

func (c *OAuthClient) GetID() string {
	return c.ClientId
}

func (c *OAuthClient) GetSecret() string {
	return c.ClientSecret
}

func (c *OAuthClient) GetDomain() string {
	return "esenyun.com"
}

func (c *OAuthClient) IsPublic() bool {
	return true
}

func (c *OAuthClient) GetUserID() string {
	return ""
}

/**
 * OAuth2 校验 ClientSecret 方法
 * TODO: 数据库保存ClientSecret
 */
func (c *OAuthClient) VerifyPassword(password string) bool {
	// sha := sha1.New()
	// sha.Write([]byte(password))
	// password := fmt.Sprintf("%x", sha.Sum(nil))
	return c.ClientSecret == password
}
