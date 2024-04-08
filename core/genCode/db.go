package genCode

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var query Query

type Query interface {
	Tables() ([]Table, error)
	TableColumn(table string) ([]Column, error)
}

func initDb(source string) *gorm.DB {

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
func NewQuery(source, dbType string) Query {
	db := initDb(source)
	switch dbType {
	case "mysql":
		return &MysqlQuery{db: db}
	default:
		panic("unsupport " + dbType)
	}
}

var _ Query = &MysqlQuery{}

type MysqlQuery struct {
	db *gorm.DB
}

func (m *MysqlQuery) Tables() ([]Table, error) {
	var result []Table
	err := m.db.Raw("select table_name as table_name,table_comment from information_schema.tables where table_schema = ?", m.db.Migrator().CurrentDatabase()).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MysqlQuery) TableColumn(table string) ([]Column, error) {
	sql := `
	SELECT COLUMN_NAME        column_name,
       DATA_TYPE          data_type,
       CASE DATA_TYPE
           WHEN 'longtext' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'varchar' THEN c.CHARACTER_MAXIMUM_LENGTH
           WHEN 'double' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'decimal' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
           WHEN 'int' THEN c.NUMERIC_PRECISION
           WHEN 'bigint' THEN c.NUMERIC_PRECISION
           ELSE '' END AS data_type_long,
       COLUMN_COMMENT     column_comment,
       CASE COLUMN_KEY
           when 'PRI' then true
           else false end as is_key,
       CASE EXTRA
           when 'auto_increment' then true
           else false end as is_auto
	FROM INFORMATION_SCHEMA.COLUMNS c
	WHERE table_name = ?
	  AND table_schema = ?
	order by c.ordinal_position
	`

	var result []Column
	err := m.db.Raw(sql, table, m.db.Migrator().CurrentDatabase()).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
