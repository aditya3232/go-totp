package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"-" gorm:"column:password"`
	TOTPKey  string `json:"-" gorm:"column:totp_key"`
}

func (User) TableName() string {
	return "users"
}
