package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_app/internal/config"
	"todo_app/internal/handler"
	taskrepo "todo_app/internal/repository/task"
	userrepo "todo_app/internal/repository/user"
	"todo_app/internal/service"
	"todo_app/internal/service/auth"
	"todo_app/internal/service/task"
	jwttoken "todo_app/internal/service/token/jwt"
	"todo_app/internal/service/user"
	"todo_app/pkg/db/postgres"
	"todo_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

// Run wires all layers together (config -> database -> repositories -> services -> handlers)
// and starts the HTTP server. I need help
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	slog.SetDefault(logger.New(cfg.Logger))

	if cfg.ServerConfigs.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	pool, err := postgres.New(context.Background(), cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()

	// Repositories: the only code that knows SQL.
	userRepo := userrepo.NewUserRepo(pool)
	taskRepo := taskrepo.NewTaskRepository(pool)

	tokens := jwttoken.NewTokenManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	srvs := &service.Service{
		AuthService: auth.NewAuthService(userRepo, tokens),
		UserService: user.NewUserService(userRepo),
		TaskService: task.NewTaskService(taskRepo),
	}

	// Handlers: HTTP only. They receive services.
	h := handler.NewHandler(srvs)

	server := &http.Server{
		Addr:              ":" + cfg.ServerConfigs.Port,
		Handler:           h.InitRoutes(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("http server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down http server")

	// Give running requests up to shutdownTimeout to finish.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	log.Println("server stopped")
	return nil

}
