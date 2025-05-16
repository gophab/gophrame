package domain

import (
	"time"

	"github.com/gophab/gophrame/domain"
)

type Task struct {
	domain.Entity
	domain.PropertiesEnabled
	Type         string     `gorm:"column:type;default:" json:"type"`                             // 任务类型：AUTO - 自动 | MANUAL - 手动 | ASYNC - 异步
	Name         string     `gorm:"column:name" json:"name"`                                      // 任务名称
	Description  *string    `gorm:"column:description;default:null" json:"description,omitempty"` // 任务描述
	Progress     float32    `gorm:"column:progress;default:0" json:"progress"`                    // 任务进度
	Status       int        `gorm:"column:status;default:0" json:"status"`                        // 任务状态: 0 - 初始化 1 - 进行中 2 - 完成 3 - 中断
	Mode         string     `gorm:"column:mode;default:void" json:"mode"`                         // 任务结果: void - 无输入 | text - 文本 | file - 文件
	Result       *string    `gorm:"column:result;default:null" json:"result,omitempty"`
	Remark       *string    `gorm:"column:remark;default:null" json:"resmark,omitempty"`
	CreatedTime  time.Time  `gorm:"column:created_time;autoCreateTime;<-:create" json:"createdTime"`
	CreatedBy    string     `gorm:"column:created_by;<-:create" json:"createdBy"`
	UpdatedTime  *time.Time `gorm:"column:updated_time;autoUpdateTime;<-:update" json:"updatedTime,omitempty"`
	FinishedTime *time.Time `gorm:"column:finished_time" json:"finishedTime,omitempty"`
	DelFlag      bool       `gorm:"column:del_flag;default:0" json:"delFlag"`
}

func (*Task) TableName() string {
	return "sys_task"
}
