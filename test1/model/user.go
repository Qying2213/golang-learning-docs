package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// TODO: 任务1 - 补充字段的 GORM 标签和 JSON 标签
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"column:username" json:"username"` // TODO: 添加 gorm 和 json 标签
	Password  string         `gorm:"column:password" json:"-"`// TODO: 添加 gorm 标签，json 标签设为 "-" 不返回密码
	Email     string         `gorm:"column:email" json:"email"`// TODO: 添加 gorm 和 json 标签
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
