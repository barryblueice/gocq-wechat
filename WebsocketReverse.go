package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

// 反向连接总事件触发

func WebsocketReverseInit(connURL string, SelfID string, self *openwechat.Self) {

	for {
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(connURL, http.Header{"X-Self-ID": []string{SelfID}})
		if err != nil {
			log.Printf("无法连接到 WebSocket 服务器: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("成功连接到 WebSocket 服务器")

		for {
			_, msgJSON, err := conn.ReadMessage()
			if err != nil {
				log.Printf("读取消息失败: %v", err)
				break
			}

			SendHandle(self, conn, msgJSON)
		}

		conn.Close()
	}
}
