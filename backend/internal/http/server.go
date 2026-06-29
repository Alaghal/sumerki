package httpserver

import (
	"database/sql"
	"errors"
	"net/http"

	"sumerki/backend/internal/http/handlers"
	appmiddleware "sumerki/backend/internal/http/middleware"
	"sumerki/backend/internal/repository"
	"sumerki/backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(database *sql.DB, jwtSecret string) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = errorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(appmiddleware.LocalCORS()))

	health := handlers.NewHealth(database)
	e.GET("/health", health.Health)
	e.GET("/ready", health.Ready)

	users := repository.NewUserRepository(database)
	auth := service.NewAuthService(users, jwtSecret)
	authHandler := handlers.NewAuthHandler(auth)
	meHandler := handlers.NewMeHandler(auth)
	kingdoms := repository.NewKingdomRepository(database)
	rulers := repository.NewRulerRepository(database)
	resources := repository.NewResourcesRepository(database)
	buildings := repository.NewBuildingRepository(database)
	rulerService := service.NewRulerService(kingdoms, rulers)
	resourcesService := service.NewResourcesService(kingdoms, resources)
	buildingService := service.NewBuildingService(kingdoms, buildings, resourcesService)
	resourcesService.SetProductionProvider(buildingService)
	kingdomService := service.NewKingdomService(kingdoms, rulerService, resourcesService, buildingService)
	kingdomHandler := handlers.NewKingdomHandler(kingdomService)
	rulerHandler := handlers.NewRulerHandler(rulerService)
	resourcesHandler := handlers.NewResourcesHandler(resourcesService)
	buildingHandler := handlers.NewBuildingHandler(buildingService)

	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)
	e.GET("/api/me", meHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/kingdoms", kingdomHandler.Create, appmiddleware.Auth(auth))
	e.GET("/api/kingdoms/me", kingdomHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/ruler/me", rulerHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/resources/me", resourcesHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/buildings/me", buildingHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/buildings/:type/upgrade", buildingHandler.Upgrade, appmiddleware.Auth(auth))

	return e
}

func errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError
	code := "internal_error"
	message := "Internal server error"

	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		status = httpErr.Code
		message = http.StatusText(status)
		if message == "" {
			message = "HTTP error"
		}
		code = "http_error"
	}

	if writeErr := handlers.JSONError(c, status, code, message); writeErr != nil {
		c.Logger().Error(writeErr)
	}
}
