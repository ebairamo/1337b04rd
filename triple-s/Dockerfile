FROM golang:1.22-alpine

WORKDIR /app
COPY . .
RUN go build -o triple-s .

EXPOSE 9000

# Используйте shell форму CMD вместо exec формы
CMD ./triple-s --port 9000 --dir /data