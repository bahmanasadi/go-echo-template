package main

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	sqlcdb "goechotemplate/api/db/model"
	accounthandler "goechotemplate/api/internal/account/handler"
	accountrepo "goechotemplate/api/internal/account/repository"
	accountservice "goechotemplate/api/internal/account/service"
	authhandler "goechotemplate/api/internal/auth/handler"
	"goechotemplate/api/internal/auth/model"
	authrepo "goechotemplate/api/internal/auth/repository"
	authservice "goechotemplate/api/internal/auth/service"
	cfg "goechotemplate/api/internal/config"
	"goechotemplate/api/pkg/validators"
	"log"
)

func main() {
	// Load configuration
	cfg, err := cfg.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pgxCfg, err := pgx.ParseConfig(cfg.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	db := stdlib.OpenDB(*pgxCfg)
	defer db.Close()

	// Initialize repository
	personRepo := accountrepo.NewPersonRepository(db)
	authRepo := authrepo.NewAuthRepository(sqlcdb.New(db))

	// Initialize service
	authService := authservice.NewAuthService(authRepo)
	personService := accountservice.NewPersonService(personRepo)

	// Initialize Echo
	e := echo.New()
	e.Validator = validators.NewCustomValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// prometheus
	e.Use(echoprometheus.NewMiddleware("goechotemplates")) // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler())         // adds route to serve gathered metrics

	// Setup routes
	personHandler := accounthandler.NewPersonHandler(authService, personService)

	authHandler := authhandler.NewAuthHandler(authService)
	e.POST("/login", authHandler.Login)

	accountGroup := e.Group("/accounts")
	accountGroup.Use(
		echojwt.WithConfig(model.DefaultJWTConfig),
	)
	accountGroup.GET("/:id", personHandler.GetPerson)
	accountGroup.POST("/", personHandler.CreatePerson)

	// Start server
	log.Fatal(e.Start(":" + cfg.Port))
}
