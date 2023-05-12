FROM golang:1.17

WORKDIR /app

COPY . .

RUN go mod init app
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/cors
RUN go get github.com/gorilla/websocket
RUN go get github.com/go-redis/redis/v8
RUN go build
EXPOSE 8080

CMD ["./app"]