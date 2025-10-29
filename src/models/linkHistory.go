package models

import "time"

// 链接表结构
type Link struct {
	ID         int       `xorm:"'link_id'         pk autoincr"`
	SourceUrl  string    `xorm:"'source_url'      notnull"`
	ShortUrl   string    `xorm:"'short_url'       index notnull unique"`
	UserID     int       `xorm:"'user_id'         notnull"`
	ClickCount int       `xorm:"'click_count'     default(0)"`
	CreatedAt  time.Time `xorm:"'created_at'      created"`
	ExpireAt   time.Time `xorm:"'expire_at'       notnull"`
}

func (l *Link) TableName() string {
	return "link"
}

// 检查链接是否过期
func (l *Link) IsExpired() bool {
	return time.Now().After(l.ExpireAt)
}
