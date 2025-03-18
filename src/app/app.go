package app

import (
	"session-restrict/configs"
	"session-restrict/src/controller"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"session-restrict/src/lib/web"
	"session-restrict/src/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	WebServer *fiber.App

	SrvAuth    *service.Auth
	SrvSession *service.Session
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	configs.LoadEnv()
	logger.InitLogger()

	a.setupDatabases()
	a.setupServices()
	a.setupHttp()

	a.WebServer.Listen(":4000")
}

func (a *App) setupDatabases() {
	err := database.ConnectPostgresSQL()
	if err != nil {
		logger.Log.Panic(err, `failed to connect to PostgreSQL`)
	}
	logger.Log.Info(`Connected to PostgreSQL`)

	err = database.ConnectRedis()
	if err != nil {
		logger.Log.Panic(err, `failed to connect to Redis`)
	}
	logger.Log.Info(`Connected to Redis`)
}

func (a *App) setupServices() {
	a.SrvAuth = service.NewAuth()
	a.SrvSession = service.NewSession()
}

func (a *App) setupHttp() {
	webServer := web.NewWebserver()
	middleware := web.NewMiddleware(webServer)
	middleware.Init()

	controller.NewPages(webServer)
	controller.NewAuth(webServer, a.SrvAuth)
	controller.NewNotification(webServer)
	controller.NewSession(webServer, a.SrvSession)

	webServer.Static(`/assets`, `./assets`, fiber.Static{
		Download:      false,
		CacheDuration: 20 * time.Second,
		Browse:        false,
	})

	a.WebServer = webServer
}
