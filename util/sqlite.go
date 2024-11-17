package util

import "gorm.io/gorm"

// GetSQLiteTables
//
//	@param db
//	@return []string
func GetSQLiteTables(db *gorm.DB) []string {
	list := []map[string]interface{}{}
	retData := []string{}
	db.Table("sqlite_master").Where(map[string]interface{}{
		"type": "table",
	}).Find(&list)
	for _, value := range list {
		for key, value := range value {
			if key == "name" {
				retData = append(retData, value.(string))
			}
		}
	}
	return retData
}

// ExecSQLiteSQL
//
//	@param db
//	@param sql
//	@return error
func ExecSQLiteSQL(db *gorm.DB, sql string) error {
	return db.Exec(sql).Error
}
