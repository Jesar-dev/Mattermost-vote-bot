package bot

import (
    "fmt"
    "net/http"
    "strings"
    "strconv"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    command := r.FormValue("text")

    parts := strings.Fields(command)
    if len(parts) == 0 {
        fmt.Fprintln(w, "Введите команду. Например: /vote create Вопрос | Вариант1 | Вариант2")
        return
    }

    switch parts[0] {
    case "create":
        pollText := strings.TrimPrefix(command, "create")
        response, err := CreatePoll(pollText)
        if err != nil {
            fmt.Fprintln(w, "Ошибка:", err)
            return
        }
        fmt.Fprintln(w, response)

    case "results":
        if len(parts) != 2 {
            fmt.Fprintln(w, "Формат: /vote results <ID>")
            return
        }

        pollID, err := strconv.Atoi(parts[1])
        if err != nil {
            fmt.Fprintln(w, "Некорректный ID голосования")
            return
        }

        response, err := GetPollResults(pollID)
        if err != nil {
            fmt.Fprintln(w, "Ошибка:", err)
            return
        }

        fmt.Fprintln(w, response)

    default:
        if len(parts) == 2 {
            pollID, err1 := strconv.Atoi(parts[0])
            option, err2 := strconv.Atoi(parts[1])
            if err1 != nil || err2 != nil {
                fmt.Fprintln(w, "Формат: /vote <ID> <номер_варианта>")
                return
            }

            response, err := Vote(pollID, option)
            if err != nil {
                fmt.Fprintln(w, "Ошибка:", err)
                return
            }

            fmt.Fprintln(w, response)
        } else {
            fmt.Fprintln(w, "Неизвестная команда. Пример: /vote 1 2")
        }
    case "close":
    if len(parts) != 2 {
        fmt.Fprintln(w, "Формат: /vote close <ID>")
        return
    }

    pollID, err := strconv.Atoi(parts[1])
    if err != nil {
        fmt.Fprintln(w, "Некорректный ID голосования")
        return
    }

    response, err := ClosePoll(pollID)
    if err != nil {
        fmt.Fprintln(w, "Ошибка:", err)
        return
    }

    fmt.Fprintln(w, response)
    case "delete":
    if len(parts) != 2 {
        fmt.Fprintln(w, "Формат: /vote delete <ID>")
        return
    }

    pollID, err := strconv.Atoi(parts[1])
    if err != nil {
        fmt.Fprintln(w, "Некорректный ID голосования")
        return
    }

    response, err := DeletePoll(pollID)
    if err != nil {
        fmt.Fprintln(w, "Ошибка:", err)
        return
    }

    fmt.Fprintln(w, response)
    }
}
