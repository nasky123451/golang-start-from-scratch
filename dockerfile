# github.com/air-verse/air@latest go version required >= 1.23

FROM golang:1.23

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD [ "air", "-c", ".air.toml" ]