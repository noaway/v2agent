//  Created by paincompiler on 8/16/16

package godao

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// Default default database name
const Default = "default"

// for multi database
var databases = map[string]*gorm.DB{}
var mu sync.Mutex

func initDatabase(c DBConfig, name ...string) error {
	dbName := Default
	if len(name) > 0 {
		dbName = name[0]
	}

	mu.Lock()
	defer mu.Unlock()

	if db, ok := databases[dbName]; ok {
		db.Close()
	}
	db, err := Open(c)
	if err != nil {
		return err
	}
	databases[dbName] = db
	return nil
}

// Open open database
func Open(c DBConfig) (*gorm.DB, error) {
	if db, err := gorm.Open(c.GetDriver(), c.GetDSN()); err != nil {
		return nil, err
	} else {
		db.SingularTable(true)
		db.LogMode(c.GetShowSQL())
		db.DB().SetConnMaxLifetime(5 * time.Minute)
		if c.GetMaxOpenConnection() > 0 {
			db.DB().SetMaxOpenConns(c.GetMaxOpenConnection())
		}
		if c.GetMaxIdleConnection() > 0 {
			db.DB().SetMaxIdleConns(c.GetMaxIdleConnection())
		}
		db.SetLogger(ormLogger{})
		return db, nil
	}
}

// Open open database
//
//    eg: postgres://1:123@localhost:10001/aaa?sslmode=disable&ShowSQL=true&MaxOpenConn=1&MaxIdleConn=1
func OpenFromURL(u string) (*gorm.DB, error) {
	config, err := parseURL(u)
	if err != nil {
		return nil, err
	}

	if db, err := gorm.Open(config.Dialect, config.URL); err != nil {
		return nil, err
	} else {
		db.SingularTable(true)
		db.LogMode(config.ShowSQL)
		db.DB().SetConnMaxLifetime(5 * time.Minute)
		if config.MaxOpenConn > 0 {
			db.DB().SetMaxOpenConns(config.MaxOpenConn)
		}
		if config.MaxIdleConn > 0 {
			db.DB().SetMaxIdleConns(config.MaxIdleConn)
		}
		db.SetLogger(ormLogger{})
		return db, nil
	}
}

func parseURL(u string) (*Config, error) {
	config := &Config{}
	if strings.HasPrefix(u, "postgres://") || strings.HasPrefix(u, "postgresql://") {
		config.Dialect = "postgres"
	} else {
		return nil, fmt.Errorf("only support postgres for now")
	}

	uu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	q := uu.Query()
	config.ShowSQL, _ = strconv.ParseBool(q.Get("ShowSQL"))
	config.MaxOpenConn, _ = strconv.Atoi(q.Get("MaxOpenConn"))
	config.MaxIdleConn, _ = strconv.Atoi(q.Get("MaxIdleConn"))
	q.Del("ShowSQL")
	q.Del("MaxOpenConn")
	q.Del("MaxIdleConn")

	uu.RawQuery = q.Encode()

	config.URL = uu.String()
	return config, nil
}

// Get get database from name(default is default), if database not exist return nil
func Get(name ...string) *gorm.DB {
	dbName := Default
	if len(name) > 0 {
		dbName = name[0]
	}
	return databases[dbName]
}
