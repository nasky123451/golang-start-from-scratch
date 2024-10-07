package websocket

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func WebsocketClient() {
	// 连接 WebSocket 服务器的 URL
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	fmt.Printf("Connecting to %s...\n", u.String())

	// 发起连接
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	// 使用 goroutine 来处理从服务器接收到的消息
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// 检查是否是关闭消息
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("connection closed normally")
					return
				}
				log.Println("read error:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
		}
	}()

	// 定时发送消息到服务器
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for i := 0; i < 5; i++ { // 发送 5 条消息后停止
		select {
		case t := <-ticker.C:
			msg := fmt.Sprintf("Hello from client at %s", t.Format("15:04:05"))
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write error:", err)
				return
			}
			fmt.Printf("Sent: %s\n", msg)
		}
	}

	// 发送结束后主动关闭连接
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("close write error:", err)
		return
	}
	log.Println("Sent close message to server, waiting to close connection...")
	time.Sleep(time.Second)
}
