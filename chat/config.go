package chat

import (
	"log"
	"os"
)

var (
	JWTSecretKey string
	RedisHost    string
	RedisPort    string
	PostgresHost string
	PostgresPort string
	PostgresDB   string
	PostgresUser string
	PostgresPass string
)

func LoadConfig() {
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JWTSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required")
	}

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	if RedisHost == "" || RedisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PORT environment variables are required")
	}

	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresDB = os.Getenv("POSTGRES_DB")
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
	if PostgresHost == "" || PostgresPort == "" || PostgresDB == "" || PostgresUser == "" || PostgresPass == "" {
		log.Fatal("PostgreSQL configuration environment variables are required")
	}
}
