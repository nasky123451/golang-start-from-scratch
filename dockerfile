# github.com/air-verse/air@latest go version required >= 1.23
# 使用 Go 的官方鏡像
FROM golang:1.23

# 設定工作目錄
WORKDIR /app

# 安裝 Air
RUN go install github.com/air-verse/air@latest

# 複製 go.mod 和 go.sum 文件
COPY go.* ./

# 下載模塊
RUN go mod download

# 複製其餘的應用代碼
COPY . .

# 檢查是否安裝 Prometheus，若未安裝則執行安裝
RUN if ! command -v prometheus &> /dev/null; then \
      apt-get update && \
      apt-get install -y prometheus; \
    fi

# 開放 8080 和 9090 端口
EXPOSE 8080
EXPOSE 9090
EXPOSE 5432

# 設定執行命令
# ENTRYPOINT [ "air", "-c", ".air.toml" ]
ENTRYPOINT [ "go", "run", "./main.go" ]
