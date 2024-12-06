package genCode

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var query Query

type Query interface {
	Tables() ([]Table, error)
	GetTableInfo(table string) (Table, error)
	TableColumn(table string) ([]Column, error)
	GetDictList() ([]DictInfo, error)
}

func initDb(source string) *gorm.DB {

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}
func NewQuery(source, dbType string) Query {
	db := initDb(source)
	switch dbType {
	case "mysql":
		query = &MysqlQuery{db: db}
		return query
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

func (m *MysqlQuery) GetTableInfo(table string) (Table, error) {
	var result Table
	err := m.db.Raw("select table_name as table_name,table_comment from information_schema.tables where table_schema = ? and table_name = ?", m.db.Migrator().CurrentDatabase(), table).Scan(&result).Error
	if err != nil {
		return Table{}, err
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

func (m *MysqlQuery) GetDictList() ([]DictInfo, error) {
	var out = make([]DictInfo, 0)
	err := m.db.Raw("select name,type from sys_dictionaries where deleted_at IS NULL and status=1").Scan(&out).Error
	return out, err
}
