package models

import (
	"time"

	"github.com/noaway/godao"
)

// Base struct
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// InitPostgre func
func InitPostgre(config godao.PostgreSQLConfig) error {
	if err := godao.InitORM(config); err != nil {
		return err
	}
	return godao.Engine.AutoMigrate(
		new(Order),
	).Error
}

type Order struct {
	Email      string    `gorm:"column:email;primary_key"`
	Package    int       `gorm:"column:package"`
	Config     string    `gorm:"column:config"`
	Price      int64     `gorm:"gorm:"column:price"` // RMB
	ExpireTime time.Time `gorm:"column:expire_time"`
	Base
}

func (Order) TableName() string {
	return "order"
}

func (o *Order) Create() error {
	return godao.Engine.Create(o).Error
}

func (o *Order) SaveConfig() error {
	return godao.Engine.Model(o).Updates(map[string]interface{}{
		"config": o.Config,
	}).Error
}

func OrderIsNotExists(email string) bool {
	return godao.IsRecordNotFound(godao.Engine.Find(&Order{}, "email = ?", email).Error)
}

func SyncOrder(order *Order) error {
	if OrderIsNotExists(order.Email) {
		return order.Create()
	}
	return order.SaveConfig()
}
