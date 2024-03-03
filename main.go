package main

import (
	"dreampicai/handler"
	"dreampicai/pkg/sb"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	// router.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(FS))))
	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fileServer))

	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return sb.Init()
}
