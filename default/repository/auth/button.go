package auth

import (
	"strings"

	"github.com/gophab/gophrame/core/inject"
	"github.com/gophab/gophrame/core/logger"
	"github.com/gophab/gophrame/core/query"

	"github.com/gophab/gophrame/default/domain/auth"

	"gorm.io/gorm"
)

type ButtonRepository struct {
	*gorm.DB `inject:"database"`
}

var buttonRepository *ButtonRepository = &ButtonRepository{}

func init() {
	inject.InjectValue("buttonRepository", buttonRepository)
}

func (b *ButtonRepository) GetById(id int64) (*auth.Button, error) {
	var result auth.Button
	if res := b.Model(&result).Where("id = ?", id).First(&result); res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, nil
	} else {
		return &result, nil
	}
}

// 根据关键词查询用户表的条数
func (b *ButtonRepository) getCounts(keyWords string) (counts int64) {
	sql := "SELECT count(*) AS counts FROM auth_button WHERE (en_name LIKE ? OR cn_name LIKE  ?) "
	if _ = b.Raw(sql, "%"+keyWords+"%", "%"+keyWords+"%").First(&counts); counts > 0 {
		return counts
	} else {
		return 0
	}
}

// 查询（根据关键词模糊查询）
func (b *ButtonRepository) Show(keyWords string, pageable query.Pageable) (totalCounts int64, temp []auth.Button) {
	totalCounts = b.getCounts(keyWords)
	if totalCounts > 0 {
		sql := `
			SELECT 
				a.*  
			FROM 
				auth_button a 
			WHERE 
				a.cn_name LIKE ? 
				OR a.en_name LIKE ? 
			LIMIT ?,?
		`
		if res := b.Raw(sql, "%"+keyWords+"%", "%"+keyWords+"%", pageable.GetOffset(), pageable.GetLimit()).Find(&temp); res.RowsAffected > 0 {
			return totalCounts, temp
		} else {
			return totalCounts, nil
		}
	}
	return 0, nil
}

//按钮编辑页的列表展示

func (a *ButtonRepository) getCountsByButtonName(cnName string) (count int64) {
	if res := a.Model(&auth.Button{}).Where("cn_name LIKE ?", "%"+cnName+"%").Count(&count); res.Error == nil {
		return count
	}
	return 0
}

func (a *ButtonRepository) List(cnName string, pageable query.Pageable) (counts int64, data []auth.Button) {
	counts = a.getCountsByButtonName(cnName)
	if counts > 0 {
		if err := a.Model(&auth.Button{}).
			Where("cn_name LIKE ?", "%"+cnName+"%").Offset(pageable.GetOffset()).Limit(pageable.GetLimit()).Find(&data); err.Error == nil {
			return
		}
	}

	return 0, nil
}

// 新增
func (b *ButtonRepository) InsertData(data *auth.Button) (bool, error) {
	data.AllowMethod = strings.ToUpper(data.AllowMethod)
	if res := b.Create(*data); res.Error == nil {
		return true, nil
	} else {
		logger.Error("Button 数据新增出错", res.Error.Error())
		return false, res.Error
	}
}

// 更新
func (b *ButtonRepository) UpdateData(data *auth.Button) (bool, error) {
	data.AllowMethod = strings.ToUpper(data.AllowMethod)
	// Omit 表示忽略指定字段(CreatedTime)，其他字段全量更新
	if res := b.Omit("CreatedTime").Save(*data); res.Error != nil {
		logger.Error("Button 数据修改出错", res.Error.Error())
		return false, res.Error
	} else {
		return true, nil
	}
}

// 删除
func (b *ButtonRepository) DeleteData(id int64) error {
	return b.Delete(b, id).Error
}
