package goBase

import (
	"database/sql"
	"reflect"
)

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
			//ent[names[i]] = &retf
			values = append(values, &retf)
		}

		rows.Scan(values...)
		for i := 0; i < flen; i++ {

			v := reflect.ValueOf(values[i])
			switch types[i].ScanType().Kind() {
			// case reflect.Int64:
			// 	ent[names[i]] = v.Elem().Int()
			case reflect.Slice:
				ent[names[i]] = string(v.Elem().Bytes())
			case reflect.String:
				ent[names[i]] = v.Elem().String()
				// case reflect.Struct:
				// 	ent[names[i]] = v.Elem().
			}

			LogDebug("输出数据 ：name : %v , t: %v kind : %v ,value: %v",
				names[i],
				types[i].ScanType().Name(),
				types[i].ScanType().Kind(), ent[names[i]])
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
