package storage

import (
    "context"
    "log"
    "time"

    "github.com/tarantool/go-tarantool/v2"
)

var Conn *tarantool.Connection

func Connect() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    dialer := tarantool.NetDialer{
        Address: "tarantool:3301", // имя контейнера в docker-compose
        User:    "guest",
    }

    opts := tarantool.Opts{
        Timeout: 3 * time.Second,
    }

    conn, err := tarantool.Connect(ctx, dialer, opts)
    if err != nil {
        log.Fatalf("Ошибка подключения к Tarantool: %s", err)
    }

    Conn = conn
    log.Println("Подключено к Tarantool")
}