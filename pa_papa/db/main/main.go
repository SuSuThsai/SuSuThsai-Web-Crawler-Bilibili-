package main

import (
	"fmt"
	"pa_papa/db"
)

func main() {
	db.DB_init()
	db.SetZhanghaoMima(db.DB, "yamada", "cc123")
	s, err := db.GetMima(db.DB, "Yamada")
	if err != nil {
		fmt.Println("出错")
	}
	fmt.Println("密码：", s)
}
