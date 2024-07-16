package internal

import (
    "database/sql"
    "log"
    "net/http"
    "tender/internal/handlers"
)

type Application struct {
    ErrorLog    *log.Logger
    InfoLog     *log.Logger
    UserHandler *handlers.UserHandler
}

func InitializeApp(db *sql.DB, errorLog, infoLog *log.Logger) *Application {
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userHandler := handlers.NewUserHandler(userService)

    return &Application{
        ErrorLog:    errorLog,
        InfoLog:     infoLog,
        UserHandler: userHandler,
    }
}

