package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smart-todo/internal/api/handlers"
	"smart-todo/internal/api/middleware"
	"smart-todo/internal/config"
	"smart-todo/internal/repository/memory"
	"smart-todo/internal/service"
	"smart-todo/pkg/jwt"
	"strconv"
	"syscall"
	"time"
)

func main() {

	cfg := config.Load()

	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)

	userRepo := memory.NewUserRepository()
	taskRepo := memory.NewTaskRepository()

	authService := service.NewAuthService(userRepo, jwtManager)
	taskService := service.NewTaskService(taskRepo)

	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	mux.HandleFunc("POST /api/tasks", authMiddleware.Authenticate(taskHandler.CreateTask))
	mux.HandleFunc("GET /api/tasks", authMiddleware.Authenticate(taskHandler.GetTasks))
	mux.HandleFunc("GET /api/tasks/overdue", authMiddleware.Authenticate(taskHandler.GetOverdue))
	mux.HandleFunc("GET /api/tasks/stats", authMiddleware.Authenticate(taskHandler.GetStats))
	mux.HandleFunc("GET /api/tasks/{id}", authMiddleware.Authenticate(taskHandler.GetTask))
	mux.HandleFunc("PUT /api/tasks/{id}", authMiddleware.Authenticate(taskHandler.UpdateTask))
	mux.HandleFunc("DELETE /api/tasks/{id}", authMiddleware.Authenticate(taskHandler.DeleteTask))
	mux.HandleFunc("PATCH /api/tasks/{id}/done", authMiddleware.Authenticate(taskHandler.MarkAsDone))

	mux.HandleFunc("GET /api/profile", authMiddleware.Authenticate(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Write([]byte("User ID: " + strconv.Itoa(userID)))
	}))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
