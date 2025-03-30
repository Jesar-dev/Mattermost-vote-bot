package main

import (
    "log"
    "net/http"

    "mattermost-vote-bot/internal/bot"
    "mattermost-vote-bot/internal/storage"
)

func main() {
    storage.Connect()
    http.HandleFunc("/command", bot.CommandHandler)

    log.Println("Bot is running on port 8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}