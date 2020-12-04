package goBase

import (
	"fmt"
	"strconv"
)

// GetRowInt64ByKey 呼气数据返回中的int
func GetRowInt64ByKey(data map[string]interface{}, key string) int64 {
	rd := GetRowStringByKey(data, key)
	return Str2Int64(rd)
}

// GetRowStringByKey 数据返回中的字符串
func GetRowStringByKey(data map[string]interface{}, key string) string {
	if r, ok := data[key]; ok {
		return fmt.Sprint(r)
	}
	return ""
}

// Str2Int64 字符串转INT64
func Str2Int64(data string) int64 {
	r, e := strconv.ParseInt(data, 10, 64)
	if e != nil {
		return 0
	}
	return r
}
