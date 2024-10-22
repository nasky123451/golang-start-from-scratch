package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 用來升級 HTTP 連接到 WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler，處理每個 WebSocket 連接
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升級 HTTP 連接到 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升級失敗:", err)
		return
	}
	defer conn.Close()

	log.Println("新客戶端連接")

	// 循環讀取客戶端的訊息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("讀取消息失敗:", err)
			break
		}
		log.Printf("接收到消息: %s", msg)

		// 將消息回傳給客戶端（可選）
		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("回傳消息失敗:", err)
			break
		}
	}
}

func main() {
	// 設置 WebSocket 路由
	http.HandleFunc("/ws", wsHandler)

	// 啟動伺服器
	fmt.Println("WebSocket 伺服器啟動在 :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("伺服器啟動失敗:", err)
	}
}
