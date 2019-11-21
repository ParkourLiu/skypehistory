package main

import (
	"fmt"
	"net/http"
	"reflect"
)

const (
	PORT = "0.0.0.0:8889"
)

var (
	svcReflectValue reflect.Value
)

func init() {
	svcReflectValue = reflect.ValueOf(&service{}) //初始化接口等待反射
}
func main() {

	http.HandleFunc("/GetHistory", interceptor)
	fmt.Println("ListenAndServe", PORT, "...")
	if err := http.ListenAndServe(PORT, nil); err != nil {
		panic(err)
	}
}
