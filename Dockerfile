FROM golang:1.23-alpine

# Установка минимальных зависимостей
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Модули
COPY go.mod go.sum ./
RUN go mod download

# Исходники
COPY . .

# Сборка
RUN go build -o eliteboard .

EXPOSE 8081

# Запуск приложения
CMD ["./eliteboard", "--port", "8081"]