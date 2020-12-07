package goBase

import (
	"log"
	"os"
)

var (
	debugLog *log.Logger
	errorLog *log.Logger
	//errorLogF *log.Logger
	sqlLog *log.Logger
	//sqlLogF   *log.Logger
)

func init() {
	// file, err := os.OpenFile("error.log",
	// 	os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalln("Failed to open error log file:", err)
	// }
	// sfile, err := os.OpenFile("sql.log",
	// 	os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalln("Failed to open error log file:", err)
	// }
	debugLog = log.New(os.Stdout, "DEBUG:", log.Ldate)

	sqlLog = log.New(os.Stdout, "SQL:", log.Ldate)
	// sqlLogF = log.New(sfile, "SQL:", log.Ldate)

	errorLog = log.New(os.Stdout, "ERROR:", log.Ldate)
	// errorLogF = log.New(file, "ERROR:", log.Ldate)
}

// LogDebug 打印调试信息
func LogDebug(format string, args ...interface{}) {
	debugLog.Printf(format, args...)
}

// LogError 打印错误信息
func LogError(format string, args ...interface{}) {
	errorLog.Printf(format, args...)
	// errorLogF.Printf(format, args...)
}

// LogSQL 打印SQL
func LogSQL(format string, args ...interface{}) {
	sqlLog.Printf(format, args...)
	//sqlLogF.Printf(format, args...)
}
