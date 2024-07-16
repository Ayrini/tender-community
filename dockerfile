# Используем официальный образ Go в качестве базового
FROM golang:1.21-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем модульные файлы и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download
	
# Копируем остальные файлы проекта в контейнер
COPY . .

COPY config/config.yaml /app/config/config.yaml

# Компилируем приложение
RUN go build -o /app/tender cmd/tender/main.go

# Определяем команду, которая будет выполнена при запуске контейнера
CMD ["/app/tender"]

