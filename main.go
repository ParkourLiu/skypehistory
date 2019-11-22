package main

import (
	"fmt"
	"mtcomm/db/mysql"
	logger "mtcomm/log"
	"net/http"
	"reflect"
)

const (
	PORT = ":8989"
)

var (
	svcReflectValue reflect.Value
	mysqlClient     mysql.MysqlClient
	log             *logger.Logger
)

func init() {
	svcReflectValue = reflect.ValueOf(&service{}) //初始化接口等待反射
	log = logger.GetDefaultLogger()
	mysqlClient = mysql.NewMysqlClient(&mysql.MysqlInfo{
		UserName:     "root",
		Password:     "root",
		IP:           "127.0.0.1",
		Port:         "3306",
		DatabaseName: "skype",
		Logger:       logger.GetDefaultLogger(),
	})
}
func main() {
	http.HandleFunc("/GetUser", interceptor)
	http.HandleFunc("/GetUserFriend", interceptor)
	http.HandleFunc("/GetUserFriendChat", interceptor)
	http.HandleFunc("/GetFriends", interceptor)
	fmt.Println("ListenAndServe", PORT, "...")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		panic(err)
	}
}
