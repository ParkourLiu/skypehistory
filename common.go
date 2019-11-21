package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//执行post请求
func DoBytesPost(url string, post []byte) (interface{}, error) {
	body := bytes.NewReader(post)
	request, err := http.NewRequest("GET", url, body)
	if request != nil {
		if request.Body != nil {
			defer func() {
				request.Body.Close()
				request.Close = true
			}()
		}
	}
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if resp != nil {
		if resp.Body != nil {
			defer func() {
				resp.Body.Close()
				resp.Close = true
			}()
		}
	}
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reqMap, _ := jsonMap(b)
	return reqMap, err
}

func jsonMap(strByte []byte) (interface{}, error) {
	var dat interface{}
	if err := json.Unmarshal(strByte, &dat); err == nil {
		return dat, nil
	} else {
		return dat, err
	}
}