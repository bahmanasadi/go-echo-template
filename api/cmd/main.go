package main

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo-contrib/echoprometheus"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	sqlcdb "goechotemplate/api/db/model"
	authhandler "goechotemplate/api/internal/auth"
	cfg "goechotemplate/api/internal/config"
	accounthandler "goechotemplate/api/internal/person"
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
	personRepo := accounthandler.NewRepository(sqlcdb.New(db))
	authRepo := authhandler.NewRepository(sqlcdb.New(db))

	// Initialize service
	authService := authhandler.NewAuthService(authRepo)
	personService := accounthandler.NewService(personRepo)

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

	authHandler := authhandler.NewHandler(authService)
	e.POST("/login", authHandler.Login)

	accountGroup := e.Group("/accounts")
	accountGroup.Use(
		echojwt.WithConfig(authhandler.DefaultJWTConfig),
	)
	accountGroup.GET("/:id", personHandler.GetPerson)
	accountGroup.POST("/", personHandler.CreatePerson)

	// Start server
	log.Fatal(e.Start(":" + cfg.Port))
}
