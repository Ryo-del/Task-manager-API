package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type server struct {
	DB    *DB
	Redis *Redis
}

type DB struct {
	*sql.DB
}
type Redis struct{}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received", "method", r.Method, "url", r.URL.Path)
		next.ServeHTTP(w, r)
	})
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
	mux := http.NewServeMux()
	//mux.HandleFunc("/login", auth.LoginHandler)

	slog.Info("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, mux))
}
