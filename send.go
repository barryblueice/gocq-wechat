package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/eatmoreapple/openwechat"
	"github.com/gorilla/websocket"
)

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
		log.Fatalf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}
}

func SendGroupText(self *openwechat.Self, data map[string]interface{}, text string) {
	params, _ := data["params"].(map[string]interface{})
	TargetIDFloat64, _ := params["group_id"].(float64)
	TargetID := strconv.Itoa(int(TargetIDFloat64))
	log.Printf("向群组 （%s）群聊消息: %s", TargetID, text)
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
		log.Fatalf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Fatalf("发送消息失败: %v", err)
	}
}

func SendHandle(self *openwechat.Self, conn *websocket.Conn, msgJSON []byte) {
	var data map[string]interface{}
	json.Unmarshal(msgJSON, &data)
	params, _ := data["params"].(map[string]interface{})
	messageType, _ := params["message_type"].(string)
	message, _ := params["message"].([]interface{})[0].(map[string]interface{})
	textData, _ := message["data"].(map[string]interface{})
	text, _ := textData["text"].(string)

	if messageType == "private" {
		SendPrivateText(self, data, text)
	} else if messageType == "group" {
		SendGroupText(self, data, text)
	} else {
		log.Fatalf("发送消息失败: 无法解析发送对象")
	}
}
