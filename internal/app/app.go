package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hong195/web-server/config"
	"github.com/hong195/web-server/internal/controller/restapi"
	"github.com/hong195/web-server/internal/repo/persistent"
	"github.com/hong195/web-server/internal/usecase/user"
	"github.com/hong195/web-server/pkg/httpserver"
	"github.com/hong195/web-server/pkg/logger"
	"github.com/hong195/web-server/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	userRepo := persistent.NewUserRepo(pg)
	userUseCase := user.New(userRepo)

	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, l, userUseCase)

	httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
