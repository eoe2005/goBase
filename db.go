package goBase

import "database/sql"

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
		return ret
	}
	names, e2 := rows.Columns()
	if e2 != nil {
		return ret
	}
	flen := len(names)
	for rows.Next() {
		values := make([]interface{}, flen)
		pvalues := make([]interface{}, flen)
		for i := 0; i < flen; i++ {
			pvalues[i] = &values[i]
		}
		rows.Scan(pvalues)
		ent := make(map[string]interface{}, 0)
		for i := 0; i < flen; i++ {
			ent[names[i]] = values[i]
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
