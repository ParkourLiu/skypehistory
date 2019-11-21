package main

import (
	"fmt"
	"time"
)

type Service interface {
	GetHistory(data map[string]string) *ReturnData
}

type service struct{}

type ReturnData struct {
	SkyId string `json:"skyId"`
}

func (s service) GetHistory(data map[string]string) *ReturnData {
	returnData := &ReturnData{}
	code, ok1 := data["code"] //1 获取好友列表，2获好友聊天记录
	head, ok2 := data["head"]
	if !ok1 || !ok2 {
		return returnData
	}
	nowUnix := fmt.Sprint(time.Now().Unix() * 1000)
	if code == "1" {
		url := "https://client-s.gateway.messenger.live.com/v1/threads/19%3Ad33cc0307370420ca2ceaf8da571b5ba%40thread.skype/consumptionhorizons"
	} else if code == "2" {

	}
	return returnData
}
