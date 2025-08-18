FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
# 下载
RUN go mod download

COPY . .

RUN go build -o main

EXPOSE 3000

CMD [ "/app/main" ]
