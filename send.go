package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/gorilla/websocket"
)

// 返回定义

type PrivateSendReturnStatus struct {
	Status  string      `json:"status"`
	Retcode int         `json:"retcode"`
	Data    bool        `json:"data"`
	Message string      `json:"message"`
	Echo    interface{} `json:"echo"`
}

type GroupSendReturnStatus struct {
	Status  string      `json:"status"`
	Retcode int         `json:"retcode"`
	Data    bool        `json:"data"`
	Message string      `json:"message"`
	Echo    interface{} `json:"echo"`
}

// 不同类型的发送

func SendPrivateText(self *openwechat.Self, data map[string]interface{}, text string) {
	params, _ := data["params"].(map[string]interface{})
	TargetIDFloat64, _ := params["user_id"].(float64)
	TargetID := strconv.Itoa(int(TargetIDFloat64))
	Friends, _ := self.Friends()
	TargetPrivateFriend := Friends.SearchByID(TargetID).First()
	log.Printf("向用户 %s（%s）发送私聊消息: %s", TargetPrivateFriend.NickName, TargetID, text)
	self.SendTextToFriend(TargetPrivateFriend, text)

	response := PrivateSendReturnStatus{
		Status:  "ok",
		Retcode: 0,
		Data:    true,
		Echo:    data["echo"],
	}
	msgJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func SendGroupText(self *openwechat.Self, data map[string]interface{}, text string) {
	params, _ := data["params"].(map[string]interface{})
	TargetIDFloat64, _ := params["group_id"].(float64)
	TargetID := strconv.Itoa(int(TargetIDFloat64))
	log.Printf("向群组 （%s）发送群聊消息: %s", TargetID, text)
	Groups, _ := self.Groups()
	TargetGroup := Groups.SearchByID(TargetID).First()
	self.SendTextToGroup(TargetGroup, text)

	response := GroupSendReturnStatus{
		Status:  "ok",
		Retcode: 0,
		Data:    true,
		Echo:    data["echo"],
	}
	msgJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func SendPrivateLocalImg(self *openwechat.Self, data map[string]interface{}, LocalFile string) {
	params, _ := data["params"].(map[string]interface{})
	TargetIDFloat64, _ := params["user_id"].(float64)
	TargetID := strconv.Itoa(int(TargetIDFloat64))
	Friends, _ := self.Friends()
	TargetPrivateFriend := Friends.SearchByID(TargetID).First()
	LocalFile = strings.ReplaceAll(LocalFile, "file:///", "")
	log.Printf("向用户 %s（%s）发送私聊本地图片: %s", TargetPrivateFriend.NickName, TargetID, LocalFile)
	img, _ := os.Open(LocalFile)
	defer img.Close()
	self.SendImageToFriend(TargetPrivateFriend, img)

	response := PrivateSendReturnStatus{
		Status:  "ok",
		Retcode: 0,
		Data:    true,
		Echo:    data["echo"],
	}
	msgJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func SendGroupLocalImg(self *openwechat.Self, data map[string]interface{}, LocalFile string) {
	params, _ := data["params"].(map[string]interface{})
	TargetIDFloat64, _ := params["group_id"].(float64)
	TargetID := strconv.Itoa(int(TargetIDFloat64))
	Groups, _ := self.Groups()
	TargetGroup := Groups.SearchByID(TargetID).First()
	LocalFile = strings.ReplaceAll(LocalFile, "file:///", "")
	log.Printf("向群组 （%s）发送群聊本地图片: %s", TargetID, LocalFile)
	img, _ := os.Open(LocalFile)
	defer img.Close()
	self.SendImageToGroup(TargetGroup, img)

	response := GroupSendReturnStatus{
		Status:  "ok",
		Retcode: 0,
		Data:    true,
		Echo:    data["echo"],
	}
	msgJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

// 发送总处理

func SendHandle(self *openwechat.Self, conn *websocket.Conn, msgJSON []byte) {
	var data map[string]interface{}
	json.Unmarshal(msgJSON, &data)
	params, _ := data["params"].(map[string]interface{})
	SendingType, _ := params["message_type"].(string)
	messages, _ := params["message"].([]interface{})

	// fmt.Println(data)

	for _, msg := range messages {
		message := msg.(map[string]interface{})
		Data, _ := message["data"].(map[string]interface{})
		text, _ := Data["text"].(string)
		LocalFile, _ := Data["file"].(string)
		messageType, _ := message["type"].(string)

		if SendingType == "private" {
			if messageType == "text" {
				SendPrivateText(self, data, text)
			} else if messageType == "image" {
				SendPrivateLocalImg(self, data, LocalFile)
			} else {
				log.Println("发送消息失败: 无法解析发送对象")
			}
		} else if SendingType == "group" {
			if messageType == "text" {
				SendGroupText(self, data, text)
			} else if messageType == "image" {
				SendGroupLocalImg(self, data, LocalFile)
			} else {
				log.Println("发送消息失败: 无法解析发送对象")
			}
		} else {
			log.Println("发送消息失败: 无法解析发送对象")
		}
	}
}
