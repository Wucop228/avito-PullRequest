package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/Wucop228/avito-PullRequest/internal/config"
	httpdelivery "github.com/Wucop228/avito-PullRequest/internal/delivery/http"
	"github.com/Wucop228/avito-PullRequest/internal/service"
)

type App struct {
	cfg  *config.Config
	db   *sql.DB
	echo *echo.Echo
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := newDB(cfg.DB)
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	teamSvc := service.NewTeamService(db)
	userSvc := service.NewUserService(db)
	prSvc := service.NewPullRequestService(db)

	teamHandler := httpdelivery.NewTeamHandler(teamSvc)
	userHandler := httpdelivery.NewUserHandler(userSvc)
	prHandler := httpdelivery.NewPullRequestHandler(prSvc)

	e.POST("/team/add", teamHandler.TeamAdd)
	e.GET("/team/get", teamHandler.TeamGet)

	e.POST("/users/setIsActive", userHandler.SetIsActive)
	e.GET("/users/getReview", prHandler.GetUserReviews)

	e.POST("/pullRequest/create", prHandler.Create)
	e.POST("/pullRequest/merge", prHandler.Merge)
	e.POST("/pullRequest/reassign", prHandler.Reassign)

	return &App{
		cfg:  cfg,
		db:   db,
		echo: e,
	}, nil
}

func newDB(cfg config.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("connected to database")
	return db, nil
}

func (a *App) RunHTTP() error {
	addr := ":" + a.cfg.Server.Port
	log.Printf("starting HTTP server on %s", addr)
	return a.echo.Start(addr)
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.echo.Shutdown(ctx); err != nil {
		return err
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return err
		}
	}

	return nil
}
