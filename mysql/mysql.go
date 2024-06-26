package mysqlUtil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func init() {
	var err error
	// 连接MySQL数据库
	DB, err = sql.Open("mysql", "root:R00t@#123!@tcp(42.193.225.126:35711)/A5")
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 尝试与数据库建立连接
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
}
