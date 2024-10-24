package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"example.com/m/chat/config"
	"example.com/m/chat/middlewares"
	"example.com/m/chat/utils"
)

var jwtKey = []byte("secret_key")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 注册用户
func RegisterUser(c *gin.Context) {

	if config.PgConn == nil {
		config.Logger.Error("PostgreSQL connection is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not available"})
		return
	}
	if config.Ctx == nil {
		config.Logger.Error("Context is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Context is not available"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		config.Logger.WithFields(logrus.Fields{
			"username": user.Username,
			"error":    err.Error(),
		}).Error("Registration failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		config.Logger.Error("Error hashing password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = config.PgConn.Exec(config.Ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, hash)
	if err != nil {
		config.Logger.WithField("username", user.Username).Error("Error registering user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
		return
	}

	config.Logger.WithField("username", user.Username).Info("User registered successfully")
	config.RegisterUserCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{"status": "User registered"})
}

// 登录用户并生成 JWT
func LoginUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		config.Logger.WithField("error", err.Error()).Error("Login failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedHash string
	err := config.PgConn.QueryRow(config.Ctx, "SELECT password FROM users WHERE username=$1", user.Username).Scan(&storedHash)
	if err != nil {
		config.Logger.Error("Invalid username or password")
		config.LoginCounter.WithLabelValues("failure").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)); err != nil {
		config.Logger.Error("Invalid username or password")
		config.LoginCounter.WithLabelValues("failure").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 JWT token
	token, err := middlewares.GenerateJWT(user.Username)
	if err != nil {
		config.Logger.Error("Error generating token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	config.LoginCounter.WithLabelValues("success").Inc()
	c.JSON(http.StatusOK, gin.H{"token": token}) // 返回 token 给前端
}

// 处理用户登出
func LogoutUser(c *gin.Context) {
	// 从请求的 JWT 中提取用户信息
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, err := middlewares.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	username := claims.Username

	// 更新用户在线状态到 Redis
	if err := utils.UpdateUserOnlineStatus(config.RedisClient, config.Ctx, username, false); err != nil {
		log.Println("Error updating online status in Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
