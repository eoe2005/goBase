package goBase

import "time"

// DBTable 表配置信息
type DBTable struct {
	DbConfName string
	TableName  string
	PK         string
}

// ORM 数据表中的字段对应
type ORM struct {
	ID       int64     `dbkey:"id"`
	CreateAt time.Time `dbkey:"create_at"`
	CreateIP string    `dbkey:"create_ip"`
	UpdateAt time.Time `dbkey:"update_at"`
	UpdateIP string    `dbkey:"update_ip"`
}

func (a *DBTable) getCon() {
	if a.PK == nil {

	}
	//return a.PK
}

// Find 查询一行
func (a *DBTable) Find(id int64) (interface{}, error) {}

// UpdateByIDAndMap 使用主键和map更新记录
func (a *DBTable) UpdateByIDAndMap(id int64, data map[string]interface{}) int {
	return 0
}

// UpdateByIDAndStruct 使用主键和结构体更新一行
func (a *DBTable) UpdateByIDAndStruct(id, int64, data interface{}) int {
	return 0
}

// Insert 使用Map插入一行
func (a *DBTable) Insert(data map[string]interface{}) int {
	return 0
}

// Save 用结构体插入一行记录
func (a *DBTable) Save(data interface{}) int {
	return 0
}

// Delete 根据主键删除一行记录
func (a *DBTable) Delete(id int64) int {
	return 0
}

//QueryRow 插叙数据
func (a *DBTable) QueryRow(whereform ...string) {
	a.QueryAll(whereform...)
}

// QueryAll 查询全部的咧
func (a *DBTable) QueryAll(whereform ...string) {}
