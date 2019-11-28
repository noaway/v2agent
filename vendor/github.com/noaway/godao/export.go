package godao

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Engine global connection
var Engine *gorm.DB // todo remove this

// InitOrm init the postgres
func InitORM(c PostgreSQLConfig) error {
	if err := initDatabase(c, Default); err != nil {
		return err
	}
	Engine = Get(Default)
	if Engine == nil {
		return fmt.Errorf("can not get database %s", Default)
	}
	return nil
}

// InitTestORM init test orm
func InitTestORM() {
	if err := initDatabase(MemorySqliteConfig{}); err != nil {
		panic(err)
	}
	Engine = Get(Default)
	if Engine == nil {
		panic(fmt.Sprintf("can not get database %s", Default))
	}
}
