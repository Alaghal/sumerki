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
	army := repository.NewArmyRepository(database)
	missions := repository.NewMissionRepository(database)
	reports := repository.NewReportRepository(database)
	patrons := repository.NewPatronRepository(database)
	raids := repository.NewRaidRepository(database)
	patronPressure := repository.NewPatronPressureRepository(database)
	rulerService := service.NewRulerService(kingdoms, rulers)
	resourcesService := service.NewResourcesService(kingdoms, resources)
	buildingService := service.NewBuildingService(kingdoms, buildings, resourcesService)
	armyService := service.NewArmyService(kingdoms, army, resourcesService, buildingService)
	missionService := service.NewMissionService(kingdoms, missions, reports, armyService, resourcesService)
	patronService := service.NewPatronService(kingdoms, patrons)
	patronPressureService := service.NewPatronPressureService(kingdoms, patrons, patronPressure, resourcesService)
	raidService := service.NewRaidService(kingdoms, raids, reports, armyService, resourcesService, buildingService)
	resourcesService.SetProductionProvider(buildingService)
	missionService.SetRaidResolver(raidService)
	patronService.SetPressureLifecycle(patronPressureService)
	kingdomService := service.NewKingdomService(kingdoms, rulerService, resourcesService, buildingService, armyService)
	kingdomHandler := handlers.NewKingdomHandler(kingdomService)
	rulerHandler := handlers.NewRulerHandler(rulerService)
	resourcesHandler := handlers.NewResourcesHandler(resourcesService)
	buildingHandler := handlers.NewBuildingHandler(buildingService)
	armyHandler := handlers.NewArmyHandler(armyService)
	missionHandler := handlers.NewMissionHandler(missionService)
	reportHandler := handlers.NewReportHandler(missionService)
	patronHandler := handlers.NewPatronHandler(patronService)
	patronPressureHandler := handlers.NewPatronPressureHandler(patronPressureService)
	raidHandler := handlers.NewRaidHandler(raidService)

	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)
	e.GET("/api/me", meHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/kingdoms", kingdomHandler.Create, appmiddleware.Auth(auth))
	e.GET("/api/kingdoms/me", kingdomHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/ruler/me", rulerHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/resources/me", resourcesHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/buildings/me", buildingHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/buildings/:type/upgrade", buildingHandler.Upgrade, appmiddleware.Auth(auth))
	e.GET("/api/army/me", armyHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/army/train", armyHandler.Train, appmiddleware.Auth(auth))
	e.GET("/api/missions/available", missionHandler.Available, appmiddleware.Auth(auth))
	e.GET("/api/missions/me", missionHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/missions/start", missionHandler.Start, appmiddleware.Auth(auth))
	e.GET("/api/reports/me", reportHandler.Me, appmiddleware.Auth(auth))
	e.GET("/api/reports/:id", reportHandler.Detail, appmiddleware.Auth(auth))
	e.POST("/api/reports/:id/read", reportHandler.MarkRead, appmiddleware.Auth(auth))
	e.GET("/api/patron/options", patronHandler.Options, appmiddleware.Auth(auth))
	e.GET("/api/patron/me", patronHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/patron/join", patronHandler.Join, appmiddleware.Auth(auth))
	e.POST("/api/patron/break", patronHandler.Break, appmiddleware.Auth(auth))
	e.GET("/api/patron/pressure", patronPressureHandler.Current, appmiddleware.Auth(auth))
	e.POST("/api/patron/pay-tribute", patronPressureHandler.PayTribute, appmiddleware.Auth(auth))
	e.POST("/api/patron/crisis-choice", patronPressureHandler.ChooseCrisis, appmiddleware.Auth(auth))
	e.GET("/api/neighbors", raidHandler.Neighbors, appmiddleware.Auth(auth))
	e.GET("/api/raids/me", raidHandler.Me, appmiddleware.Auth(auth))
	e.POST("/api/raids/start", raidHandler.Start, appmiddleware.Auth(auth))

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
