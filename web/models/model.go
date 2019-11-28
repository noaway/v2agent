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
		new(User),
		new(Commodity),
		new(Order),
		new(VPSNode),
		new(ProxyAccount),
	).Error
}

type User struct {
	ID    int    `gorm:"column:id"`
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
	Role  int    `gorm:"column:role"`

	Base
}

func (User) TableName() string { return "user" }

type Commodity struct {
	ID          int       `gorm:"column:id"`
	ExpireTime  time.Time `gorm:"column:expire_time"`
	Price       int       `gorm:"column:price"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`

	Base
}

func (Commodity) TableName() string { return "commodity" }

type Order struct {
	ID          int `gorm:"column:id"`
	UserID      int `gorm:"column:user_id"`
	CommodityID int `gorm:"column:commodity_id"`

	Base
}

func (Order) TableName() string { return "order" }

type VPSNode struct {
	ID            int    `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	User          string `gorm:"column:user"`
	Host          string `gorm:"column:host"`
	PrivateKey    string `gorm:"column:private_key"`
	Region        string `gorm:"column:region"`
	StartScript   string `gorm:"column:start_script"`
	TroubleScript string `gorm:"column:trouble_script"`

	Base
}

func (VPSNode) TableName() string { return "vps_node" }

type ProxyAccount struct {
	UUID           string `gorm:"column:uuid"`
	Name           string `gorm:"column:name"`
	Type           string `gorm:"column:type"`
	Server         string `gorm:"column:server"`
	Port           int    `gorm:"column:port"`
	AlterId        string `gorm:"column:alter_id"`
	Cipher         string `gorm:"column:cipher"`
	Network        string `gorm:"column:network"`
	WsPath         string `gorm:"column:ws_path"`
	TLS            bool   `gorm:"column:tls"`
	SkipCertVerify bool   `gorm:"column:skip_cert_verify"`

	Base
}

func (ProxyAccount) TableName() string { return "proxy_account" }
