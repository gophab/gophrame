package domain

type CountryArea struct {
	Code     string `gorm:"column:code" json:"code"`
	AreaCode int64  `gorm:"column:area_code" json:"areaCode"`
	Name     string `gorm:"column:name" json:"name"`
}

func (*CountryArea) TableName() string {
	return "sys_country_area"
}
