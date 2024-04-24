package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/gorilla/websocket"
)

type PrivateMessage struct {
	Time        int64  `json:"time"`
	SelfID      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageID   int32  `json:"message_id"`
	UserID      int64  `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      Sender `json:"sender"`
}

type GroupMessage struct {
	Time        int64  `json:"time"`
	SelfID      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageID   int32  `json:"message_id"`
	GroupID     int64  `json:"group_id"`
	UserID      int64  `json:"user_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      Sender `json:"sender"`
}

type Sender struct {
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
}

func get_timestamp() int {
	currentTime := time.Now()
	timestampMillis := currentTime.UnixNano() / int64(time.Millisecond)
	timestampString := strconv.FormatInt(timestampMillis, 10)
	timestampInt, err := strconv.Atoi(timestampString)
	if err != nil {
		panic(err)
	}
	return timestampInt
}

func RecievePrivateText(self *openwechat.Self, MsgID int32, SenderID string, message string) {
	timestamp := get_timestamp()
	TargetSenderID, _ := strconv.ParseInt(SenderID, 10, 64)
	msg := PrivateMessage{
		Time:        int64(timestamp),
		SelfID:      int64(self.ID()),
		PostType:    "message",
		MessageType: "private",
		SubType:     "friend",
		UserID:      TargetSenderID,
		Message:     message,
		RawMessage:  message,
		Font:        14,
		Sender: Sender{
			Nickname: self.NickName,
		},
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}

func RecieveGroupText(self *openwechat.Self, MsgID int32, SenderID string, message string, GroupID string, IsAt bool) {
	timestamp := get_timestamp()
	TargetSenderID, _ := strconv.ParseInt(SenderID, 10, 64)
	TargetGroupID, _ := strconv.ParseInt(GroupID, 10, 64)

	var CQMessage string

	message = strings.ReplaceAll(message, fmt.Sprintf("@%s", self.NickName), "")

	if IsAt {
		CQMessage = fmt.Sprintf("[CQ:at,qq=%d]%s", self.ID(), message)
	} else {
		CQMessage = message
	}

	msg := GroupMessage{
		Time:        int64(timestamp),
		SelfID:      int64(self.ID()),
		PostType:    "message",
		MessageType: "group",
		SubType:     "normal",
		GroupID:     TargetGroupID,
		UserID:      TargetSenderID,
		Message:     CQMessage,
		RawMessage:  message,
		Font:        14,
		Sender: Sender{
			Nickname: self.NickName,
		},
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("无法序列化消息: %v", err)
	}

	err = conn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}
}
