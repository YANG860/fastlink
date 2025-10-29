package models

import "time"

// 用户表结构
type UserHistory struct {
	ID        int       `xorm:"'user_id'     pk  autoincr"`
	RemovedAt time.Time `xorm:"'removed_at'  created"`
	Account   string    `xorm:"'account'     notnull"`
	PwHash    string    `xorm:"'pw_hash'     notnull"`

	LinkCount       int       `xorm:"'link_count'      default(0)"`
	LatestCreatedAt time.Time `xorm:"'latest_created_at'"`
}

func (u *UserHistory) TableName() string {
	return "user_history"
}
