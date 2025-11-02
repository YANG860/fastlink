package db

import "time"	
// 链接表结构
type Link struct {
	ID         int       `xorm:"'link_id'         pk autoincr"           redis:"link_id"`
	SourceUrl  string    `xorm:"'source_url'      "                      redis:"source_url"`
	ShortUrl   string    `xorm:"'short_url'       index notnull unique"  redis:"short_url"`
	UserID     int       `xorm:"'user_id'         notnull"               redis:"user_id"`
	ClickCount int       `xorm:"'click_count'     default(0)"            redis:"click_count"`
	CreatedAt  time.Time `xorm:"'created_at'      created"               redis:"created_at"`
	ExpireAt   time.Time `xorm:"'expire_at'       notnull"               redis:"expire_at"`
}

func (l *Link) TableName() string {
	return "link"
}


// 检查链接是否过期
func (l *Link) IsExpired() bool {
	return time.Now().After(l.ExpireAt)
}