package chat

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	redisClient *redis.Client
	pgConn      *pgxpool.Pool
	ctx         = context.Background()
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients     = make(map[*websocket.Conn]string)
	sessionTTL  = 10 * time.Minute
	mu          sync.Mutex
	logger      = logrus.New()
	authKey     = "YOUR_GENERATED_AUTH_KEY"
	secretKey   = "YOUR_GENERATED_SECRET_KEY"
	Log         *logrus.Logger

	// Prometheus metrics
	registerUserCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_registration_total",
			Help: "Total number of user registrations",
		},
		[]string{"status"},
	)

	loginCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_login_total",
			Help: "Total number of user logins",
		},
		[]string{"status"},
	)
	messageCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "chat_messages_total",
		Help: "Total number of chat messages sent",
	}, []string{"room"})
)
