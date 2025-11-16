package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/Wucop228/avito-PullRequest/internal/config"
	"github.com/Wucop228/avito-PullRequest/internal/delivery/http"
	"github.com/Wucop228/avito-PullRequest/internal/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	e := echo.New()

	teamSvc := service.NewTeamService(db)
	teamHandler := http.NewTeamHandler(teamSvc)

	userSvc := service.NewUserService(db)
	userHandler := http.NewUserHandler(userSvc)

	prSvc := service.NewPullRequestService(db)
	prHandler := http.NewPullRequestHandler(prSvc)

	e.POST("/team/add", teamHandler.TeamAdd)
	e.GET("/team/get", teamHandler.TeamGet)

	e.POST("/users/setIsActive", userHandler.SetIsActive)
	e.GET("/users/getReview", prHandler.GetUserReviews)

	e.POST("/pullRequest/create", prHandler.Create)
	e.POST("/pullRequest/merge", prHandler.Merge)
	e.POST("/pullRequest/reassign", prHandler.Reassign)

	if err := e.Start(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
