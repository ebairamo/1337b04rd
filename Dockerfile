FROM golang:1.23-alpine

# Установка зависимостей для Delve
RUN apk add --no-cache git bash libc6-compat

# Установка Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Модули
COPY go.mod go.sum ./
RUN go mod download

# Исходники
COPY . .

# Сборка с поддержкой отладки
RUN go build -gcflags="all=-N -l" -o eliteboard .

EXPOSE 8081 40000

# Запуск Delve с передачей аргументов
CMD ["dlv", "exec", "./eliteboard", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "--log", "--log-dest=stdout", "--", "--port=8081"]
