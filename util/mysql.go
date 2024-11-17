package util

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetMySQLTables
//
//	@param gormDB
//	@return []string
func GetMySQLTables(gormDB *gorm.DB) []string {
	tables := []string{}
	list := []map[string]interface{}{}
	if err := gormDB.Raw("show tables").Scan(&list).Error; err != nil {
		return tables
	}

	for _, value := range list {
		for _, value := range value {
			tables = append(tables, value.(string))
		}
	}

	return tables
}

// TableDesc
type TableDesc struct {
	Default interface{} `json:"Default"`
	Extra   string      `json:"Extra"`
	Field   string      `json:"Field"`
	Key     string      `json:"Key"`
	Null    string      `json:"Null"`
	Type    string      `json:"Type"`
}

// GetMySQLTableFields
//
//	@param gormDB
//	@param table
//	@return []TableDesc
func GetMySQLTableFields(gormDB *gorm.DB, table string) []TableDesc {
	list := []TableDesc{}
	fileds := []map[string]interface{}{}
	if err := gormDB.Raw("desc `" + table + "`").Find(&fileds).Error; err != nil {
		return list
	}
	bytes, _ := json.Marshal(fileds)
	if err := json.Unmarshal(bytes, &list); err != nil {
		fmt.Printf("json unmarshal error: %v", err)
		return list
	}
	return list
}

// GetCreateTableSQL
//
//	@param gormDB
//	@param table
//	@return string
func GetCreateTableSQL(gormDB *gorm.DB, table string) string {
	data := map[string]interface{}{}
	if err := gormDB.Raw("show create table `" + table + "`").Scan(&data).Error; err != nil {
		return ""
	}
	if value, ok := data["Create Table"]; ok {
		if stringValue, ok := value.(string); ok {
			return stringValue
		}
	}

	return ""
}

var gormDSN string = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc="

// GetGormDB
//
//	@param database
//	@return *gorm.DB
//	@return error
func GetGormDB(database string) (*gorm.DB, error) {
	user := GetMySQLUser()
	password := GetMySQLPassword()
	host := GetMySQLHost()
	port := GetMySQLPort()

	dsn := fmt.Sprintf(gormDSN, user, password, host, port, database)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// UpdateDataByID ...
//
//	@param db
//	@param table
//	@param id
//	@param field
//	@param value
//	@return error
func UpdateDataByID(db *gorm.DB, table string, id interface{}, data map[string]interface{}) error {
	return db.Table(table).Where(map[string]interface{}{"id": id}).Updates(data).Error
}

// DeleteDataByID
//
//	@param db
//	@param table
//	@param id
//	@return error
func DeleteDataByID(db *gorm.DB, table string, id interface{}) error {
	sql := fmt.Sprintf("delete from %s where id = ?", table)
	return db.Exec(sql, id).Error
}

// QueryList
//
//	@param db
//	@param table
//	@param where
//	@return []map[string]interface{}
//	@return error
func QueryList(db *gorm.DB, table string, where map[string]interface{}) ([]map[string]interface{}, error) {
	list := []map[string]interface{}{}
	if err := db.Table(table).Where(where).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// CreateData
//
//	@param db
//	@param table
//	@param data
//	@return error
func CreateData(db *gorm.DB, table string, data map[string]interface{}) error {
	return db.Table(table).Create(&data).Error
}
