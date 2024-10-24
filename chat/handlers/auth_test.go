package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/chat/config"
	"example.com/m/chat/handlers"
	"example.com/m/chat/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// 初始化測試環境
func init() {
	config.Init()
	metrics.InitMetrics()
}

func TestRegisterUser(t *testing.T) {

	// 設置 gin 引擎
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	// 模擬有效的用戶註冊請求
	validUser := `{"username": "testuser", "password": "testpassword"}`
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(validUser)))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 使用 httptest Recorder 來模擬 HTTP 回應
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 驗證回應狀態碼
	assert.Equal(t, http.StatusOK, w.Code)

	// 驗證回應 JSON
	expectedBody := `{"status":"User registered"}`
	assert.JSONEq(t, expectedBody, w.Body.String())

	// 刪除已註冊的用戶
	err = deleteUser("testuser") // 使用您定義的刪除用戶函數
	if err != nil {
		t.Fatalf("Failed to delete user: %v\n", err)
	}
}

func TestRegisterUserDatabaseError(t *testing.T) {
	// 每次測試前清空已註冊的指標
	prometheus.Unregister(metrics.RegisterUserCounter)
	prometheus.Unregister(metrics.LoginCounter)

	// 重新初始化 Prometheus 指標
	metrics.InitMetrics()
	// 設置 gin 引擎
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/register", handlers.RegisterUser)

	// 模擬當 PostgreSQL 連接不可用的情況
	config.PgConn = nil

	invalidUser := `{"username": "testuser", "password": "testpassword"}`
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte(invalidUser)))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 驗證回應狀態碼
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 驗證回應 JSON
	expectedBody := `{"error":"Database connection is not available"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

// 假設有一個函數用來刪除用戶
func deleteUser(username string) error {
	_, err := config.PgConn.Exec(config.Ctx, "DELETE FROM users WHERE username = $1", username)
	return err
}
