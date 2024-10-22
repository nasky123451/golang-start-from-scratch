package chat

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" // Redis 客戶端
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Message struct {
	Room    string    `json:"room"`
	From    string    `json:"from"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func initChat() {
	var err error
	// 初始化 Redis 客戶端
	redisClient, err = initRedis()

	// 初始化 PostgreSQL
	pgConn, err = initDB()

	if err := checkAndCreateTableChat(pgConn); err != nil {
		log.Fatalf("Error checking/creating chat table: %v", err)
	}

	// 註冊 Prometheus 指標
	prometheus.MustRegister(registerUserCounter)
	prometheus.MustRegister(loginCounter)
	prometheus.MustRegister(messageCounter)

	// 日誌配置
	logger.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		log.Fatal("Failed to log to file, using default stderr")
	}
}

func ChatServer() {
	initChat()

	r := gin.Default()

	// CSRF 保護
	//r.Use(gin.WrapH(csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))(r)))

	// 添加 CORS 支持
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                             // 替换为你的前端地址
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},                        // 确保允许 OPTIONS 方法
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "Authorization"}, // 添加您需要的自定义头
		AllowCredentials: true,
	}))

	// 处理 OPTIONS 请求
	r.OPTIONS("/register", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 允许所有源
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Status(http.StatusNoContent) // 返回 204 No Content
	})

	r.OPTIONS("/online-users", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Prometheus 指标
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.GET("/online-users", getOnlineUsers)
	r.GET("/chat-history", getChatHistory)
	r.GET("/ws", handleWebSocket) // WebSocket 处理路由

	// JWT 认证中间件
	r.Use(middlewareJWT())

	r.Run(":8080")
}

// 处理 WebSocket 连接时更新在线用户状态
func handleWebSocket(c *gin.Context) {
	// 升级 HTTP 连接到 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	// 等待接收身份验证消息
	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		// 处理身份验证消息
		if msg["type"] == "auth" {
			tokenString := msg["token"]
			log.Printf("Received token: %s", tokenString)

			claims, err := ParseToken(tokenString)

			if err == nil {
				username := claims.Username
				clients[conn] = username // 将用户添加到连接列表
				log.Printf("User %s connected", username)
				broadcastUserStatus(username, true) // 广播用户上线状态

				// 更新用户在线状态到 Redis
				if err := updateUserOnlineStatus(username, true); err != nil {
					log.Println("Error updating online status in Redis:", err)
				}
			} else {
				log.Println("Could not parse claims")
				break
			}
		}

		// 处理聊天消息
		if msg["type"] == "message" {
			room := msg["room"]
			from := msg["from"]
			content := msg["content"]
			timeStr := msg["time"]

			msgTime, err := time.Parse(time.RFC3339, timeStr)
			if err != nil {
				log.Println("Invalid message time:", err)
				continue
			}

			message := Message{
				Room:    room,
				From:    from,
				Content: content,
				Time:    msgTime,
			}

			if err := saveMessageToDB(message); err != nil {
				log.Println("Error saving message to DB:", err)
				continue
			}

			broadcastMessageToRoom(room, message)
		}
	}

	// 处理用户断开连接
	username := clients[conn]
	delete(clients, conn)
	log.Printf("User %s disconnected", username)

	// 更新用户在线状态到 Redis
	if err := updateUserOnlineStatus(username, false); err != nil {
		log.Println("Error updating online status in Redis:", err)
	}

	// 广播用户下线消息
	broadcastUserStatus(username, false)
}

// 使用 Redis 存储和获取在线用户
func updateUserOnlineStatus(username string, online bool) error {
	if online {
		// 用户上线，设置键值对并设置过期时间为 1 小时
		return SetKey(redisClient, ctx, username, "online", time.Hour)
	} else {
		// 用户下线，删除键
		return DeleteKey(redisClient, ctx, username)
	}
}

func saveMessageToDB(message Message) error {
	_, err := pgConn.Exec(ctx, "INSERT INTO chat_messages (room, sender, content, time) VALUES ($1, $2, $3, $4)",
		message.Room, message.From, message.Content, message.Time)
	return err
}

func broadcastMessageToRoom(room string, message Message) {
	for client, _ := range clients {
		err := client.WriteJSON(gin.H{
			"type":    "message",
			"room":    message.Room,
			"from":    message.From,
			"content": message.Content,
			"time":    message.Time,
		})
		if err != nil {
			log.Println("Error broadcasting message:", err)
			client.Close() // 如果广播失败则关闭连接
			delete(clients, client)
		}
	}
}

// 广播用户状态
func broadcastUserStatus(username string, online bool) {
	status := "offline"
	if online {
		status = "online"
	}
	for client := range clients {
		err := client.WriteJSON(gin.H{"type": "userStatus", "username": username, "status": status})
		if err != nil {
			log.Println("Error broadcasting user status:", err)
		}
	}
}

func getOnlineUsers(c *gin.Context) {
	// 获取在线用户列表
	onlineUsers := []string{}
	keys, err := redisClient.Keys(ctx, "*").Result() // 获取匹配特定模式的所有键
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching online users"})
		return
	}

	for _, key := range keys {
		value, err := redisClient.Get(ctx, key).Result() // 獲取鍵的值
		if err != nil {
			log.Printf("Error fetching value for key %s: %v", key, err)
			continue
		}
		if value == "online" { // 檢查值是否為 "online"
			onlineUsers = append(onlineUsers, key)
		}
	}

	log.Printf("Current online users: %v", onlineUsers)
	c.JSON(http.StatusOK, gin.H{"onlineUsers": onlineUsers})
}

func getChatHistory(c *gin.Context) {
	room := c.Query("room")
	date := c.Query("date") // 格式为 YYYY-MM-DD

	// 解析日期
	var startDate time.Time
	var endDate time.Time
	var err error

	if room == "" {
		room = "general"
	}

	if date == "" {
		// 如果没有提供日期，则使用当前日期
		startDate = time.Now().Truncate(24 * time.Hour)
		endDate = startDate.Add(24 * time.Hour)
	} else {
		startDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		endDate = startDate.Add(24 * time.Hour)
	}

	// 查询聊天记录
	rows, err := pgConn.Query(ctx, "SELECT sender, content, time FROM chat_messages WHERE room = $1 AND time >= $2 AND time < $3 ORDER BY time ASC", room, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching chat history"})
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.From, &msg.Content, &msg.Time); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning message"})
			return
		}
		msg.Room = room
		messages = append(messages, msg)
	}

	// 如果没有找到消息，则返回一个状态和消息
	if len(messages) == 0 {
		c.JSON(http.StatusOK, gin.H{"messages": []Message{}, "status": "No messages found for the selected date."})
		return
	}

	// 返回找到的消息
	c.JSON(http.StatusOK, gin.H{"messages": messages, "status": "Success"})
}

func registerUser(c *gin.Context) {
	if pgConn == nil {
		logger.Error("PostgreSQL connection is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not available"})
		return
	}
	if ctx == nil {
		logger.Error("Context is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context is not available"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.WithFields(logrus.Fields{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Registration failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = pgConn.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hash)
	if err != nil {
		logger.WithField("username", user.Username).Error("Error registering user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	logger.WithField("username", user.Username).Info("User registered successfully")
	registerUserCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{"status": "User registered"})
}

func loginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.WithField("error", err.Error()).Error("Login failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedHash string
	err := pgConn.QueryRow(ctx, "SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedHash)
	if err != nil {
		logger.Error("Invalid username or password")
		loginCounter.WithLabelValues("failure").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)); err != nil {
		logger.Error("Invalid username or password")
		loginCounter.WithLabelValues("failure").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 JWT token
	token, err := generateJWT(user.Username)
	if err != nil {
		logger.Error("Error generating token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	loginCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{"token": token}) // 返回 token 给前端
}

// 生成 JWT token 的函数
func generateJWT(username string) (string, error) {
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 设置 token 过期时间为72小时
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key")) // 确保将密钥替换为您的安全密钥
}

func ParseToken(tokenString string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err == nil && tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func middlewareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 解析和验证 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 检查 token 签名方法是否正确
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 返回验证签名所需的密钥
			return []byte("your-secret-key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 继续处理请求
		c.Next()
	}
}
