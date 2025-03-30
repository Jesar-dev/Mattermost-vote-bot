# Mattermost Vote Bot

Бот для проведения голосований внутри чатов мессенджера **Mattermost**.
Реализован на языке **Go** с использованием **Tarantool** для хранения данных и запускается в контейнерах через **Docker Compose**.

---

## Функциональность

- Создание голосования
- Голосование за вариант ответа
- Просмотр результатов голосования
- Завершение голосования
- Удаление голосования

---

## Быстрый старт

### 1. Клонируйте репозиторий
```bash
git clone https://github.com/Jesar-dev/mattermost-vote-bot.git
cd mattermost-vote-bot
```

### 2. Запустите контейнеры
```bash
docker-compose up --build
```

### 3. Отправляйте команды на локальный сервер
Примеры через `curl`:

#### Создать голосование
```bash
curl -X POST http://localhost:8080/command -d "text=create Язык программирования? | Go | Rust | Python"
```

#### Проголосовать
```bash
curl -X POST http://localhost:8080/command -d "text=1 2"
```

#### Посмотреть результаты
```bash
curl -X POST http://localhost:8080/command -d "text=results 1"
```

#### Завершить голосование
```bash
curl -X POST http://localhost:8080/command -d "text=close 1"
```

#### Удалить голосование
```bash
curl -X POST http://localhost:8080/command -d "text=delete 1"
```

---

## Используемые технологии
- Go 1.21
- Tarantool 2.x
- Docker & Docker Compose

---

## Структура проекта
```
.
├── cmd/                # Точка входа (main.go)
├── internal/
│   ├── bot/            # Логика бота (polls, команды)
│   └── storage/        # Подключение к Tarantool
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## Обратная связь
Автор: Геннадий Торлак  
Контакты: torlak.g@phystech.edu
