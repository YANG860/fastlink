package db

import "time"

// 用户表结构
type User struct {
	ID           int       `xorm:"'user_id'        pk  autoincr"`
	RegisteredAt time.Time `xorm:"'created_at'     created"`
	Account      string    `xorm:"'account'        notnull"`
	PwHash       string    `xorm:"'pw_hash'        notnull"`
	Valid        bool      `xorm:"'valid'          default(1) notnull"`
	TokenVersion int       `xorm:"'token_version'  default(0) notnull"`

	LinkCount       int       `xorm:"'link_count'     default(0)"`
	LatestCreatedAt time.Time `xorm:"'latest_created_at'"`
}

func (u *User) TableName() string {
	return "user"
}

// 检查用户是否有效
func (u *User) IsValid() bool {
	return u.Valid
}
