package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	//skyId:="19%3Ad33cc0307370420ca2ceaf8da571b5ba%40thread.skype"
	//url := "https://client-s.gateway.messenger.live.com/v1/users/ME/conversations/"+skyId+"/messages?view=supportsExtendedHistory%7Cmsnp24Equivalent%7CsupportsMessageProperties&pageSize=200&startTime=1573115226830"//1573115226830
	//head := map[string]string{
	//	"RegistrationToken": " registrationToken=U2lnbmF0dXJlOjI6Mjg6QVFRQUFBQXh2dm9HYXpMR1hNcWt1OEQ0OCtpWDtWZXJzaW9uOjY6MToxO0lzc3VlVGltZTo0OjE5OjUyNDg3ODU5MjE0OTA0NjU5ODA7RXAuSWRUeXBlOjc6MTo4O0VwLklkOjI6MjY6bGl2ZTouY2lkLjJmNmJhOWQxYWZlODBiMjA7RXAuRXBpZDo1OjM2OjMyOThjZDAwLWZmZmYtZmZmZi0wOTI3LTlhZGU0YjIyNTVmNTtFcC5Mb2dpblRpbWU6NzoxOjA7RXAuQXV0aFRpbWU6NDoxOTo1MjQ4Nzg1OTIxNDgzMjc4NjgxO0VwLkF1dGhUeXBlOjc6MjoxNTtFcC5FeHBUaW1lOjQ6MTk6NTI0ODc4Njc4NTQzNzM4NzkwNDtVc3IuTmV0TWFzazoxMToxOjI7VXNyLlhmckNudDo2OjE6MDtVc3IuUmRyY3RGbGc6MjowOjtVc3IuRXhwSWQ6OToxOjA7VXNyLkV4cElkTGFzdExvZzo0OjE6MDtVc2VyLkF0aEN0eHQ6Mjo0NDQ6Q2xOcmVYQmxWRzlyWlc0YWJHbDJaVG91WTJsa0xqSm1ObUpoT1dReFlXWmxPREJpTWpBQkExVnBZeFF4THpFdk1EQXdNU0F4TWpvd01Eb3dNQ0JCVFF4T2IzUlRjR1ZqYVdacFpXUWdDK2l2MGFsckx3QUFBQUFBQUVBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBYWJHbDJaVG91WTJsa0xqSm1ObUpoT1dReFlXWmxPREJpTWpBQUFBQUFBQUFBQUFBSFRtOVRZMjl5WlFBQUFBQUVBQUFBQUFBQUFBQUFBQUFnQytpdjBhbHJMd0FBQUFBQUFBQUFBQUFBQUFBQUFBQUFBUnBzYVhabE9pNWphV1F1TW1ZMlltRTVaREZoWm1VNE1HSXlNQUFBQUFBQW5WYlhYUWNBQUFBSVNXUmxiblJwZEhrT1NXUmxiblJwZEhsVmNHUmhkR1VJUTI5dWRHRmpkSE1PUTI5dWRHRmpkSE5WY0dSaGRHVUlRMjl0YldWeVkyVU5RMjl0YlhWdWFXTmhkR2x2YmhWRGIyMXRkVzVwWTJGMGFXOXVVbVZoWkU5dWJIa0FBQT09Ow==; expires=1574479901; endpointId={3298cd00-ffff-ffff-0927-9ade4b2255f5}",
	//}
	////DoBytesPost("GET", url, head, nil)
	//
	//url2:="https://client-s.gateway.messenger.live.com/v1/users/ME/conversations?view=supportsExtendedHistory%7Cmsnp24Equivalent&pageSize=2&startTime=1&targetType=Passport%7CSkype%7CLync%7CThread%7CAgent%7CShortCircuit%7CPSTN%7CSmsMms%7CFlxt%7CNotificationStream%7CCast%7CCortanaBot%7CModernBots%7CsecureThreads%7CInviteFree"
	//DoBytesPost("GET", url2, head, nil)
	a := " registrationToken=U2lnbmF0dXJlOjI6Mj"
	fmt.Println(strings.Index(a, "registrationToken="))
	a = a[1+len("registrationToken="):]
	fmt.Println(a)
}

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
