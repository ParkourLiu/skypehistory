package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

//拦截器
func interceptor(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format("2006-01-02 15:04:05")
	ip := RemoteIp(r)
	if r.Method!="POST" {
		fmt.Println(now,ip,"request method is",r.Method)
		return
	}
	RequestURI := r.RequestURI
	fmt.Println(">>>", now, ip, RequestURI)
	RequestURI = strings.Replace(r.RequestURI, "/", "", 1)
	index := strings.Index(RequestURI, "?")
	if index > -1 {
		RequestURI = RequestURI[:index]
	}
	reqStr,reqMap, err := decodeRequest(r)
	if err != nil {
		fmt.Println("decodeRequest:", reqStr, err)
		return
	}

	ret1 := Apply(svcReflectValue, RequestURI, []interface{}{reqMap}) //反射调用
	jsonBytes, err := json.Marshal(ret1[0].Interface())               //json化
	if err != nil {
		fmt.Println("json.Marshal:", err)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		fmt.Println("w.Write:", err)
		return
	}
}

//反射调用
func Apply(value_f reflect.Value, methodName string, args []interface{}) []reflect.Value {
	method := value_f.MethodByName(methodName)
	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}
	return method.Call(in)
}

//格式化参数
func decodeRequest(r *http.Request) (string,map[string]string, error) {
	defer r.Body.Close()
	reqBytes,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		return "",nil,err
	}
	reqMap,err:=json2Map(reqBytes)
	return string(reqBytes),reqMap,err
}

//获取ip
func RemoteIp(r *http.Request) string {
	//now := time.Now().Format("2006-01-02 15:04:05")
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}
func json2Map(strByte []byte) (map[string]string, error) {
	var dat map[string]string
	if err := json.Unmarshal(strByte, &dat); err == nil {
		return dat, nil
	} else {
		return dat, err
	}
}