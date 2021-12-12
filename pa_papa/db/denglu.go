package db

import (
	"database/sql"
	"errors"
	"fmt"
)

func GetMima(DB *sql.DB, zhanghao string) (string, error) {
	if len(zhanghao) == 0 {
		return "", errors.New("账号不能为空")
	}
	sqlstr := `select user_pass from wp_users where user_login=?`
	s := ""
	err := DB.QueryRow(sqlstr, zhanghao).Scan(&s)
	if err != nil {
		fmt.Println("找不到该账号密码：", err)
		return "", err
	}
	fmt.Println("查找密码成功")
	return s, nil
}

//0,插入成功
//1，账号或者密码不符合标准
//2，账号已存在
//3，插入错误
func SetZhanghaoMima(DB *sql.DB, zhanghao, mima string) (int, error) {
	if len(zhanghao) == 0 {
		return 1, errors.New("账号不能为空")
	}
	if len(mima) == 0 {
		return 1, errors.New("密码不能为空")
	}
	_, err := GetMima(DB, zhanghao)
	if err == nil {
		fmt.Println("账号存在")
		return 2, err
	}
	sqlstr := "insert into wp_users(wp_login,user_pass) values (?,?)"
	ret, err := DB.Exec(sqlstr, &zhanghao, &mima)
	if err != nil {
		fmt.Println("插入账号密码错误：", err)
		return 3, err
	}
	fmt.Println("成功插入账号密码")
	theID, _ := ret.LastInsertId()
	fmt.Println("插入id：", theID)
	return 0, nil
}
