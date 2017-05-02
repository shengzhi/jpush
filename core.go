// Package jpush 实现了极光推送RestAPI第三版
package jpush

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Platform 推送平台
type Platform byte

// 推送平台定义
const (
	IOS Platform = 1 << iota
	Android
	WinPhone
	AllPlatform Platform = IOS | Android | WinPhone
)

type M map[string]interface{}

// PushRequest 推送请求
type PushRequest struct {
	Platform     Platform      `json:"platform"`
	Audience     Audience      `json:"audience"`
	Notification *Notification `json:"notification"`
	Message      *Message      `json:"message"`
	Options      Option        `json:"options"`
}

// PushReplay 推送结果
type PushReplay struct {
	SendNo string `json:"sendno"`
	MsgID  string `json:"msg_id"`
}

// JPushClient 极光推送客户端
type JPushClient struct {
	appkey, secret string
	basicauth      string
	client         *http.Client
}

// NewJPushClient 创建极光推送客户端实例
// 如果queueLen 大于0， 则启用异步推送
func NewJPushClient(key, secret string, queueLen int) *JPushClient {
	client := JPushClient{appkey: key, secret: secret}
	client.basicauth = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", key, secret)))
	client.client = &http.Client{Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return &client
}

const apigateway = "https://api.jpush.cn"

// Push 执行消息推送
func (j *JPushClient) Push(req PushRequest) (PushReplay, error) {
	return j.call(req, fmt.Sprintf("%s/v3/push", apigateway))
}

// PushToAll 向所有用户推送消息通知
func (j *JPushClient) PushToAll(content string) (PushReplay, error) {
	req := PushRequest{
		Platform:     AllPlatform,
		Audience:     AudienceAll(),
		Notification: &Notification{Alert: content},
	}
	return j.Push(req)
}

// PushToUser 推送消息给指定的用户
func (j *JPushClient) PushToUser(content string, ids ...string) (PushReplay, error) {
	aud := &Audience{}
	aud = aud.AddObject(AudienceRegID(ids...))
	req := PushRequest{
		Platform:     AllPlatform,
		Audience:     *aud,
		Notification: &Notification{Alert: content},
	}
	return j.Push(req)
}

// PushToAnyTag 推送消息给任意tag
func (j *JPushClient) PushToAnyTag(content string, tag ...string) (PushReplay, error) {
	aud := &Audience{}
	aud = aud.AddObject(AudienceTag(tag...))
	req := PushRequest{
		Platform:     AllPlatform,
		Audience:     *aud,
		Notification: &Notification{Alert: content},
	}
	return j.Push(req)
}

const validateapigateway = "https://api.jpush.cn/v3/push/validate"

// Validate 消息校验
func (j JPushClient) Validate(req PushRequest) (PushReplay, error) {
	return j.call(req, validateapigateway)
}

func (j JPushClient) call(r PushRequest, url string) (PushReplay, error) {
	var replay PushReplay
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(r)
	if err != nil {
		return replay, fmt.Errorf("序列化错误:%v", err)
	}
	request, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return replay, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", j.basicauth))
	response, err := j.client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return replay, err
	}

	if response.StatusCode != 200 {
		data, _ := ioutil.ReadAll(response.Body)
		return replay, fmt.Errorf("发送失败:%s", string(data))
	}
	err = json.NewDecoder(response.Body).Decode(&replay)
	if err != nil {
		return replay, fmt.Errorf("反序列化错误:%v", err)
	}
	return replay, nil
}
