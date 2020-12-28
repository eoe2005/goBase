package goBase

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"
)

// DBTable 表结构
type DBTable struct {
	table String
	con   String
	db    *sql.DB
}

// InitDBTable 初始化配置
func InitDBTable(tableName, conName string, db *sql.DB) *DBTable {
	return &DBTable{
		table: tableName,
		con:   conName,
		db:    db,
	}
}

// DBGetRow 获取一行记录
func DBGetRow(db *sql.DB, format string, args ...interface{}) map[string]interface{} {
	data := DBGetAll(db, format, args...)
	if len(data) > 0 {
		return data[0]
	}
	return nil
}

// DBGetAll 获取全部的记录
func DBGetAll(r *sql.DB, format string, args ...interface{}) []map[string]interface{} {
	rows, e := r.Query(format, args...)
	defer rows.Close()
	ret := make([]map[string]interface{}, 0)
	if e != nil {
		LogError("SQL 错误 %v", e)
		return ret
	}
	types, e3 := rows.ColumnTypes()
	if e3 != nil {
		LogError("SQL 查询结果格式错误 %v", e3)
		return ret
	}
	names, e2 := rows.Columns()
	if e2 != nil {
		LogError("SQL 错误 %v", e2)
		return ret
	}
	flen := len(names)

	for rows.Next() {
		ent := make(map[string]interface{}, flen)
		values := make([]interface{}, 0, flen)
		for i := 0; i < flen; i++ {
			var retf interface{}
			values = append(values, &retf)
		}

		rows.Scan(values...)
		for i := 0; i < flen; i++ {

			v := reflect.ValueOf(values[i])
			switch types[i].ScanType().Kind() {
			case reflect.Int64:
				switch reflect.ValueOf(v.Elem().Interface()).Type().Kind() {
				case reflect.Int64:
					ent[names[i]] = v.Elem().Interface().(int64)
				case reflect.Slice:
					ent[names[i]], _ = strconv.ParseInt(string(v.Elem().Interface().([]byte)), 10, 64)
				}
			case reflect.Slice:
				ent[names[i]] = string(v.Elem().Interface().([]byte))
			case reflect.String:
				ent[names[i]], _ = strconv.ParseInt(string(v.Elem().Interface().([]byte)), 10, 64)
			case reflect.Struct:
				ent[names[i]] = string(v.Elem().Interface().([]byte))
			default:
				ent[names[i]] = v.Elem().Interface()
			}

			//LogDebug("输出数据 ：name : %v , t: %v -> %v kind : %v -> %v ,value: %v",
			//	names[i],
			//	types[i].ScanType().Name(), v.Elem().Type().Name(),
			//	types[i].ScanType().Kind(),
			//	reflect.ValueOf(v.Elem().Interface()).Type().Kind(),
			//	//v.Elem().Type().Kind(),
			//	ent[names[i]])
		}
		ret = append(ret, ent)
	}
	return ret
}

// DBUpdate 更新记录
func DBUpdate(db *sql.DB, format string, args ...interface{}) int64 {
	r, e := db.Exec(format, args...)
	if e != nil {
		LogError("SQL 错误 %v", e)
		return 0
	}
	ret, e2 := r.RowsAffected()
	if e2 != nil {
		LogError("SQL 错误 %v", e2)
		return 0
	}
	return ret
}

// DBInsert 插入记录
func DBInsert(db *sql.DB, format string, args ...interface{}) int64 {
	r, e := db.Exec(format, args...)
	if e != nil {
		LogError("SQL 错误 %v", e)
		return 0
	}
	ret, e2 := r.LastInsertId()
	if e2 != nil {
		LogError("SQL 错误 %v", e2)
		return 0
	}
	return ret
}

// DBDelete 删除记录
func DBDelete(db *sql.DB, format string, args ...interface{}) int64 {
	r, e := db.Exec(format, args...)
	if e != nil {
		LogError("SQL 错误 %v", e)
		return 0
	}
	ret, e2 := r.RowsAffected()
	if e2 != nil {
		LogError("SQL 错误 %v", e2)
		return 0
	}
	return ret
}

// InsertData 插入数据
func (t *DBTable) InsertData(data map[string]interface{}) int64 {
	le := len(data)
	keys := make([]string, 0, le)
	v1 := make([]string, 0, le)
	vals := make([]interface{}, 0, len)
	for k, v := range data {
		keys = append(keys, k)
		vals = append(vals, v)
		v1 = append(v1,"?")
	}
	return DBInsert(t.db, fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", t.table,strings.Join(keys,","),strings.Join(v1,",")),...vals)
}
// UpdateDataById 根据ID更新数据
func (t *DBTable) UpdateDataById(id,int64,data map[string]interface{})int64{
	le := len(data)
	sets := make([]string, 0, le)
	vals := make([]interface{}, 0, len)
	for k, v := range data {
		sets = append(keys, fmt.Sprintf("%v=?",k))
		vals = append(vals, v)
	}

	return DBUpdate(t.db,fmt.Sprintf("UPDATE %v SET %v WHERE id=?",t.table,strings.Join(sets,",")),id)
}

// Find 根据ID查询一条记录
func (t *DBTable) Find(id interface{}){
	return DBGetRow(t.db,fmt.Sprintf("SELECT * FROM %v WHERE id=?",t.table),id)
}
// FindByWhere 根据where条件查询一行
func func (t *DBTable) FindByWhere(format string,...args interface{}){}
// DeleteByID 删除一行数据
func func (t *DBTable) DeleteByID(id int64)int64{
	return DBDelete(t.db,fmt.Sprintf("DELETE FROM %v WHERE id=?",t.table),id)
}
