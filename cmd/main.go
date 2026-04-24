package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"taskmanager/auth"
	"taskmanager/repo"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type server struct {
	userRepo *repo.UserRepository
}

func CROSHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	// Register handlers here and pass s.userRepo into them or into a service layer.
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(w, r, s.userRepo)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterHandler(w, r, s.userRepo)
	})

	return CROSHeadersMiddleware(Middleware(mux))
}

func main() {

	viper.SetConfigFile("config/config.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	port := ":" + viper.GetString("server.port")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	userRepo := repo.NewUserRepository(db)

	svr := &server{
		userRepo: userRepo,
	}

	slog.Info("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, svr.routes()))
}
