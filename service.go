package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"time"
)

//type Service interface {
//	GetHistory(reqData ReqData) *ReturnData
//}

type service struct{}

type ReqData struct {
	SkyIds    []map[string]string `json:"skyIds"`
	Head      string              `json:"Head"`
	Name      string              `json:"name"`      //别名
	Version   string              `json:"version"`   //消息唯一Id
	Friend_id string              `json:"friend_id"` //消息唯一Id
	HadeMap   map[string]string   `json:"hadeMap"`
}

type ReturnData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewReturnData() *ReturnData {
	return &ReturnData{Code: "100", Msg: "succeed"}
}

type GetFriendsStuct struct {
	Conversations []conversations `json:"conversations"`
	Metadata      metadata        `json:"_metadata"`
}

type conversations struct {
	Id string `json:"id"`
}

type AllContacts struct {
	Contacts []contacts `json:"contacts"`
}

type contacts struct {
	Person_id    string `json:"person_id"`
	Display_name string `json:"display_name"`
}

func (s service) GetFriends(reqData ReqData) *ReturnData {
	returnData := NewReturnData()
	if reqData.Head == "" || reqData.Name == "" {
		returnData.Code = "101"
		returnData.Msg = "参数错误"
		return returnData
	}
	headMap, err := headResolver(reqData.Head)
	if err != nil {
		returnData.Code = "101"
		returnData.Msg = "参数错误"
		return returnData
	}

	err = iu_name(reqData.Name)
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	//近期联系人
	getUrl := "https://client-s.gateway.messenger.live.com/v1/users/ME/conversations?view=supportsExtendedHistory%7Cmsnp24Equivalent&pageSize=200&startTime=1&targetType=Passport%7CSkype%7CLync%7CThread%7CAgent%7CShortCircuit%7CPSTN%7CSmsMms%7CFlxt%7CNotificationStream%7CCast%7CCortanaBot%7CModernBots%7CsecureThreads%7CInviteFree"
nextPage:
	reqBytes, err := DoBytesPost("GET", getUrl, headMap, nil)
	if err != nil {
		returnData.Code = "102"
		returnData.Msg = "请求近期联系人接口错误：" + err.Error()
		return returnData
	}
	getFriendsStuct, err := json2GetFriendsStuct(reqBytes)
	if err != nil {
		returnData.Code = "103"
		returnData.Msg = "请求最近联系人接口返回数据结构异常：" + err.Error()
		return returnData
	}
	for _, frind := range getFriendsStuct.Conversations {
		frindType := "0"
		if strings.HasSuffix(frind.Id, ".skype") { //群
			frindType = "1"
		}
		err = iu_friends(reqData.Name, frind.Id, "", frindType)
		if err != nil {
			returnData.Code = "104"
			returnData.Msg = err.Error()
			return returnData
		}
	}
	if getFriendsStuct.Metadata.BackwardLink != "" { //下一页
		getUrl = getFriendsStuct.Metadata.BackwardLink
		goto nextPage
	}

	//所有联系人
	getUrl = "https://edge.skype.com/pcs/contacts/v2/users/self"
	token := headMap["Authentication"]
	prefixStr := "skypetoken="
	headMap["X-SkypeToken"] = token[strings.Index(token, prefixStr)+len(prefixStr):]

	reqBytes, err = DoBytesPost("GET", getUrl, headMap, nil)
	if err != nil {
		returnData.Code = "102"
		returnData.Msg = "请求全部联系人接口错误：" + err.Error()
		return returnData
	}
	allContacts, err := json2AllContacts(reqBytes)
	if err != nil {
		returnData.Code = "103"
		returnData.Msg = "请求全部联系人接口返回数据结构异常：" + err.Error()
		return returnData
	}
	for _, frind := range allContacts.Contacts {
		frindType := "0"
		if strings.HasSuffix(frind.Person_id, ".skype") { //群
			frindType = "1"
		}
		err = iu_friends(reqData.Name, frind.Person_id, frind.Display_name, frindType)
		if err != nil {
			returnData.Code = "104"
			returnData.Msg = err.Error()
			return returnData
		}
	}

	//遍历聊天记录
	reqList, err := s_friends(reqData)
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	reqData.SkyIds = reqList
	reqData.HadeMap = headMap
	err = GetHistory(reqData)
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	return returnData
}

func (s service) GetUser(reqData ReqData) *ReturnData {
	returnData := NewReturnData()
	userList, err := s_name()
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	returnData.Data = userList
	return returnData
}

func (s service) GetUserFriend(reqData ReqData) *ReturnData {
	returnData := NewReturnData()
	if reqData.Name == "" {
		returnData.Code = "101"
		returnData.Msg = "参数错误"
		return returnData
	}
	friendList, err := s_friends(reqData)
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	returnData.Data = friendList
	return returnData
}

func (s service) GetUserFriendChat(reqData ReqData) *ReturnData {
	returnData := NewReturnData()
	if reqData.Name == "" || reqData.Friend_id == "" {
		returnData.Code = "101"
		returnData.Msg = "参数错误"
		return returnData
	}
	if reqData.Version == "" {
		reqData.Version = fmt.Sprint(time.Now().Unix() * 1000)
	}
	chatList, err := s_chat(reqData)
	if err != nil {
		returnData.Code = "104"
		returnData.Msg = err.Error()
		return returnData
	}
	returnData.Data = chatList
	return returnData
}

type GetHistoryStuct struct {
	Messages []message `json:"messages"`
	Metadata metadata  `json:"_metadata"`
}
type message struct {
	Version     string `json:"version"`     //发消息的时间戳/主键
	Messagetype string `json:"messagetype"` //消息类型
	Content     string `json:"content"`     //消息内容
	Type        string `json:"type"`        //类型
	From        string `json:"from"`        //谁发送的
}
type metadata struct {
	BackwardLink string `json:"backwardLink"`
}

func GetHistory(reqData ReqData) error {
	if reqData.SkyIds == nil || len(reqData.SkyIds) < 1 {
		return nil
	}

	for _, friend := range reqData.SkyIds {
		skypeId := friend["friend_id"]
		getUrl := "https://client-s.gateway.messenger.live.com/v1/users/ME/conversations/" + skypeId + "/messages?view=supportsExtendedHistory%7Cmsnp24Equivalent%7CsupportsMessageProperties&pageSize=200&startTime=1514736000000"
	nextPage:
		reqBytes, err := DoBytesPost("GET", getUrl, reqData.HadeMap, nil)
		if err != nil {
			log.Debug(skypeId, err.Error())
			continue
		}
		getHistoryStuct, err := json2GetHistoryStuct(reqBytes)
		if err != nil {
			log.Debug(err)
			continue
		}
		//插数据库
		for _, chat := range getHistoryStuct.Messages {
			err = iu_chat(reqData.Name, skypeId, chat)
			if err != nil {
				log.Debug(err)
				continue
			}
		}
		if getHistoryStuct.Metadata.BackwardLink != "" { //下一页
			getUrl = getHistoryStuct.Metadata.BackwardLink
			goto nextPage
		}
	}

	return nil
}

//解析头
func headResolver(headStr string) (map[string]string, error) {
	headMap := map[string]string{}
	headStr = strings.Replace(headStr, "\r", "", -1)
	heads := strings.Split(headStr, "\n")
	for _, head := range heads {
		headKV := strings.Split(head, ":")
		if len(headKV) == 2 {
			headMap[headKV[0]] = headKV[1]
		}
	}
	if _, ok := headMap["RegistrationToken"]; !ok { //没有关键token
		return nil, errors.New("无用户token")
	}
	if _, ok := headMap["Authentication"]; !ok { //没有关键token
		return nil, errors.New("无用户token")
	}
	return headMap, nil
}
