package main

import (
	"fmt"
	"log"
	"net/http"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/limiter/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	redisStore, err := store.NewRedisStore(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	limiterInstance := limiter.NewRateLimiter(cfg, redisStore)

	r := chi.NewRouter()
	r.Use(limiterInstance.Middleware) // ← Aqui usamos o método da instância

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Rate limiter OK!"))
	})

	fmt.Println("Servidor rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
