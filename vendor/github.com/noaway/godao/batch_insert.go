package godao

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type tabler interface {
	TableName() string
}

func generateSQL(db *gorm.DB, values []interface{}) string {
	var (
		tnames        = map[string][]*gorm.Scope{}
		sortTableName = make([]string, 0, len(values))
		fir           = false
		columns       []string
	)

	for _, v := range values {
		if tab, ok := v.(tabler); ok {
			tn := tab.TableName()
			if str, ok := tnames[tn]; ok {
				tnames[tn] = append(str, db.NewScope(v))
			} else {
				sortTableName = append(sortTableName, tn)
				scope := db.NewScope(v)
				tnames[tn] = []*gorm.Scope{scope}
				if !fir {
					for _, field := range scope.Fields() {
						if field.IsNormal {
							if field.Field.Kind() == reflect.Ptr && field.IsBlank {
								continue
							}
							if !field.IsPrimaryKey || (field.IsPrimaryKey && !field.IsBlank) {
								columns = append(columns, scope.Quote(field.DBName))
							}
						}
					}
					fir = true
				}
			}
		} else {
			panic("no implementation TableName function")
		}
	}

	sqlBuf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(sqlBuf)
	sqlBuf.Reset()

	sort.Strings(sortTableName)
	for _, tableName := range sortTableName {
		sqlBuf.WriteString(fmt.Sprintf("insert into %v (%v) values ", tableName, strings.Join(columns, ",")))
		scopes := tnames[tableName]
		for i := 0; i < len(scopes); i++ {
			genValuesSQL(scopes[i].Fields(), sqlBuf)
			if i < len(scopes)-1 {
				sqlBuf.WriteString(",")
			}
		}
		sqlBuf.WriteString(";")
	}
	return sqlBuf.String()
}

func genValuesSQL(fields []*gorm.Field, buf *bytes.Buffer) {
	buf.WriteString("(")
	first := true
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		if !field.IsNormal || (field.Field.Kind() == reflect.Ptr && field.IsBlank) {
			continue
		}

		if strings.ToLower(field.Name) == "id" {
			if field.IsBlank {
				continue
			}
		}

		if !first {
			buf.WriteString(",")
		} else {
			first = false
		}

		buf.WriteString(formatField(field))
	}
	buf.WriteString(")")
}

func formatField(field *gorm.Field) string {
	if field.IsBlank && (field.Name == "CreatedAt" || field.Name == "UpdatedAt") {
		return fmt.Sprintf("'%s'", time.Now().Format(time.RFC3339Nano))
	}

	if t, ok := field.Field.Interface().(time.Time); ok {
		return fmt.Sprintf("'%s'", t.Format(time.RFC3339Nano))
	}
	if t, ok := field.Field.Interface().(*time.Time); ok {
		return fmt.Sprintf("'%s'", t.Format(time.RFC3339Nano))
	}

	return ReplaceQuotes(fmt.Sprintf("%v", field.Field))
}

func ReplaceQuotes(val string) string {
	return strings.Join([]string{"'", strings.Replace(val, "'", "''", -1), "'"}, "")
}

func BatchInsert(db *gorm.DB, values []interface{}) error {
	return db.Exec(generateSQL(db, values)).Error
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func isBlank(value reflect.Value) bool {
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
