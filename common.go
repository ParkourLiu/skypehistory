package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

//执行post请求
func DoBytesPost(method string, url string, headMap map[string]string, post []byte) ([]byte, error) {
	body := bytes.NewReader(post)
	request, err := http.NewRequest(method, url, body)
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
	//request.Header.Set("Content-Type", "application/json")
	for k, v := range headMap {
		request.Header[k] = []string{v}
	}
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
	if resp.StatusCode != 200 { //请求不成功
		return nil, errors.New(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

func json2GetHistoryStuct(strByte []byte) (GetHistoryStuct, error) {
	var dat GetHistoryStuct
	if err := json.Unmarshal(strByte, &dat); err == nil {
		return dat, nil
	} else {
		return dat, err
	}
}
func json2GetFriendsStuct(strByte []byte) (GetFriendsStuct, error) {
	var dat GetFriendsStuct
	if err := json.Unmarshal(strByte, &dat); err == nil {
		return dat, nil
	} else {
		return dat, err
	}
}

func json2AllContacts(strByte []byte) (AllContacts, error) {
	var dat AllContacts
	if err := json.Unmarshal(strByte, &dat); err == nil {
		return dat, nil
	} else {
		return dat, err
	}
}