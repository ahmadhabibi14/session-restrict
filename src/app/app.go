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
	SrvAuth   *service.Auth
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
	database.ConnectPostgresSQL()
	logger.Log.Info(`Connected to PostgreSQL`)

	database.ConnectRedis()
	logger.Log.Info(`Connected to Redis`)
}

func (a *App) setupServices() {
	a.SrvAuth = service.NewAuth()
}

func (a *App) setupHttp() {
	webServer := web.NewWebserver()
	middleware := web.NewMiddleware(webServer)
	middleware.Init()

	controller.NewPages(webServer)
	controller.NewAuth(webServer, a.SrvAuth)

	webServer.Static(`/assets`, `./assets`, fiber.Static{
		Download:      false,
		CacheDuration: 20 * time.Second,
		Browse:        false,
	})

	a.WebServer = webServer
}
