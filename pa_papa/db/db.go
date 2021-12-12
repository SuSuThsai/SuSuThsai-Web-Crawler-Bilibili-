package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DB_init() (err error) {
	dsn := "root:1H9tNa,4l6D*@tcp(localhost:3306)/wordpress"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("连接数据库出错：", err)
		return err
	}
	err = DB.Ping()
	if err != nil {
		fmt.Println("ping数据库出错：", err)
		return err
	}
	fmt.Println("连接数据库成功")
	DB.SetMaxIdleConns(10)
	DB.SetMaxOpenConns(20)
	return
}
