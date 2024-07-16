package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"tender/internal/config"
	"tender/internal/handlers"
	"tender/internal/repositories"
	"tender/internal/services"

	"github.com/rs/cors"
)

func openDB(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < 15; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Error opening database: %v", err)
			return nil, err
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database")
			return db, nil
		}

		log.Printf("Error pinging database (attempt %d): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, err
}

// dsd
func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

type Application struct {
	ErrorLog    *log.Logger
	InfoLog     *log.Logger
	UserHandler *handlers.UserHandler
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/signup", app.UserHandler.SignUp)
	mux.HandleFunc("/users", app.UserHandler.GetAllUsers)
	return mux
}

func main() {
	cfg := config.LoadConfig()

	port := os.Getenv("PORT")
	if port != "" {
		port = ":" + port
	} else {
		port = ":4000"
	}

	addr := flag.String("addr", port, "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.Database.URL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}(db)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	app := &Application{
		ErrorLog:    errorLog,
		InfoLog:     infoLog,
		UserHandler: userHandler,
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:19006", "exp://192.168.1.219:8081", "exp://192.168.1.82:8081", "timetodo://"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
	})

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      addSecurityHeaders(logRequests(c.Handler(app.routes()))),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
