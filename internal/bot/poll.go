package bot

import (
    "fmt"
    "log"
    "strings"
    "sync"

    "mattermost-vote-bot/internal/storage"
    tarantool "github.com/tarantool/go-tarantool/v2"
)

type Poll struct {
    ID       int
    Question string
    Options  []string
    Votes    map[int]int // option index -> count
    Active   bool
}

var (
    polls      = make(map[int]*Poll)
    pollsMutex sync.Mutex
    nextID     = 1
)

func toInt(v interface{}) int {
    switch t := v.(type) {
    case int:
        return t
    case int8:
        return int(t)
    case int16:
        return int(t)
    case int32:
        return int(t)
    case int64:
        return int(t)
    case float64:
        return int(t)
    case uint64:
        return int(t)
    default:
        return 0
    }
}

func CreatePoll(text string) (string, error) {
    parts := strings.Split(text, "|")
    if len(parts) < 2 {
        return "", fmt.Errorf("Нужно указать вопрос и хотя бы один вариант через |")
    }

    question := strings.TrimSpace(parts[0])
    options := []string{}
    for _, opt := range parts[1:] {
        options = append(options, strings.TrimSpace(opt))
    }

    if len(options) < 2 {
        return "", fmt.Errorf("Минимум два варианта ответа")
    }

    seqResp, err := storage.Conn.Do(
        tarantool.NewEvalRequest("return box.sequence.poll_id:next()"),
    ).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка генерации ID: %v", err)
    }
    log.Printf("Тип resp[0]: %T", seqResp[0])

    pollID := uint64(seqResp[0].(int8))
    votes := make([]int, len(options)) // [0, 0, 0]

    req := tarantool.NewInsertRequest("polls").Tuple([]interface{}{
        pollID,
        question,
        options,
        votes,
        true,
    })

    _, err = storage.Conn.Do(req).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка записи в Tarantool: %v", err)
    }

    response := fmt.Sprintf("Голосование создано (ID %d):\n%s", pollID, question)
    for i, opt := range options {
        response += fmt.Sprintf("\n%d. %s", i+1, opt)
    }

    return response, nil
}

func Vote(pollID int, optionIndex int) (string, error) {
    req := tarantool.NewSelectRequest("polls").
        Index("primary").
        Limit(1).
        Iterator(tarantool.IterEq).
        Key([]interface{}{uint64(pollID)})

    data, err := storage.Conn.Do(req).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка чтения из Tarantool: %v", err)
    }
    if len(data) == 0 {
        return "", fmt.Errorf("Голосование с ID %d не найдено", pollID)
    }

    tuple := data[0].([]interface{})
    active := tuple[4].(bool)
    if !active {
        return "", fmt.Errorf("Голосование уже завершено")
    }

    votes := tuple[3].([]interface{})
    if optionIndex < 1 || optionIndex > len(votes) {
        return "", fmt.Errorf("Некорректный номер варианта")
    }

    votes[optionIndex-1] = toInt(votes[optionIndex-1]) + 1

    updateReq := tarantool.NewReplaceRequest("polls").Tuple([]interface{}{
        tuple[0],
        tuple[1],
        tuple[2],
        votes,
        tuple[4],
    })

    _, err = storage.Conn.Do(updateReq).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка обновления голосования: %v", err)
    }

    options := tuple[2].([]interface{})
    return fmt.Sprintf("Ваш голос за «%s» засчитан!", options[optionIndex-1]), nil
}

func GetPollResults(pollID int) (string, error) {
    req := tarantool.NewSelectRequest("polls").
        Index("primary").
        Iterator(tarantool.IterEq).
        Limit(1).
        Key([]interface{}{uint64(pollID)})

    data, err := storage.Conn.Do(req).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка запроса к Tarantool: %v", err)
    }
    if len(data) == 0 {
        return "", fmt.Errorf("Голосование с ID %d не найдено", pollID)
    }

    tuple := data[0].([]interface{})
    question := tuple[1].(string)
    options := tuple[2].([]interface{})
    rawVotes := tuple[3].([]interface{})
    total := 0
    votes := make([]int, len(rawVotes))
    for i, v := range rawVotes {
        votes[i] = toInt(v)
        total += votes[i]
    }

    result := fmt.Sprintf("Результаты голосования #%d:\n%s", pollID, question)
    for i := range options {
        percent := 0
        if total > 0 {
            percent = int(float64(votes[i]) / float64(total) * 100)
        }
        result += fmt.Sprintf("\n%d. %s — %d голос(ов) (%d%%)", i+1, options[i], votes[i], percent)
    }

    return result, nil
}

func ClosePoll(pollID int) (string, error) {
    req := tarantool.NewSelectRequest("polls").
        Index("primary").
        Iterator(tarantool.IterEq).
        Limit(1).
        Key([]interface{}{uint64(pollID)})

    data, err := storage.Conn.Do(req).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка чтения из Tarantool: %v", err)
    }
    if len(data) == 0 {
        return "", fmt.Errorf("Голосование с ID %d не найдено", pollID)
    }

    tuple := data[0].([]interface{})
    active := tuple[4].(bool)
    if !active {
        return "", fmt.Errorf("Голосование уже завершено")
    }

    update := tarantool.NewReplaceRequest("polls").Tuple([]interface{}{
        tuple[0],
        tuple[1],
        tuple[2],
        tuple[3],
        false,
    })

    _, err = storage.Conn.Do(update).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка обновления голосования: %v", err)
    }

    return fmt.Sprintf("Голосование #%d завершено", pollID), nil
}

func DeletePoll(pollID int) (string, error) {
    req := tarantool.NewDeleteRequest("polls").
        Index("primary").
        Key([]interface{}{uint64(pollID)})

    _, err := storage.Conn.Do(req).Get()
    if err != nil {
        return "", fmt.Errorf("Ошибка при удалении: %v", err)
    }

    return fmt.Sprintf("Голосование #%d удалено", pollID), nil
}