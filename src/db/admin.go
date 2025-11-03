package db

// 管理员表结构
type Admin struct {
	ID int `xorm:"'admin_id'        pk  autoincr"`

	Account string `xorm:"'account'        notnull index"`
	PwHash  string `xorm:"'pw_hash'        notnull"`
}

func (u *Admin) TableName() string {
	return "admin"
}


func AuthenticateAdmin(account string, pwHash string) (bool, error) {
	var admin Admin
	has, err := SQLEngine.Where("account = ? AND pw_hash = ?", account, pwHash).Get(&admin)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}
