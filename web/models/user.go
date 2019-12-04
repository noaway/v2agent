package models

type User struct {
	ID    int    `gorm:"column:id,not null;primary_key;auto_increment"`
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
	Role  int    `gorm:"column:role"`

	Base
}

func (User) TableName() string { return "user" }
