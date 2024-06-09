# стадия 1. Сборка проекта(не будет включен в финальный образ)
FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/sport-plus/main.go 

# стадия 2. Запуск в отдельном контейнере
FROM alpine:3.14

# Создаем директорию, в которой в последующем будем работать и хранить бинарник
RUN mkdir /app

# костыль доступа к storage host машины
ARG UID=1000
ARG GID=1000

RUN adduser -D -g ${GID} -u ${UID} appuser

WORKDIR /app

# Копируем файл из под прошлого этапа, для пользователя root("защита" от перезаписи запускаемого бинарника)
COPY --from=builder --chown=root:root /build/main .

# На всякий выдаем права на чтение и выполнение другим(не root)
RUN chmod 755 main

# перевести в режим продакшена
# ENV GIN_MODE=release

# понижаем права, для запуска кода(чтобы код был запущен без ЛИШНИХ привелегий)
USER appuser

CMD ["/app/main"]