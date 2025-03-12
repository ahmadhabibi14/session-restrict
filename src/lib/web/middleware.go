package web

import (
	"fmt"
	"os"
	"session-restrict/src/lib/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

type Middleware struct {
	webServer *fiber.App
}

func NewMiddleware(webServer *fiber.App) *Middleware {
	return &Middleware{
		webServer: webServer,
	}
}

func (m *Middleware) Init() {
	m.RateLimiter()
	m.Recover()
	m.Logger()
	m.ContentSecurityPolicy()
	m.RequestID()
}

func (m *Middleware) RequestID() {
	m.webServer.Use(requestid.New(requestid.Config{
		Header: fiber.HeaderXRequestID,
		Generator: func() string {
			return uuid.NewString()
		},
	}))
}

func (m *Middleware) Recover() {
	m.webServer.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, err interface{}) {
			logger.Log.Error(err, "received unexpected panic error at "+c.Path())
		},
	}))
}

func (m *Middleware) RateLimiter() {
	m.webServer.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 2 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			message := "you have exceeded your rate limit, please try again a few moments later"
			switch c.Method() {
			case fiber.MethodGet:
				return c.Render("error", fiber.Map{
					`Title`:       fmt.Sprintf("%d - %s", fiber.StatusTooManyRequests, `Too Many Requests`),
					`Description`: message,
				})
			default:
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					`errors`: message,
				})
			}
		},
	}))
}

func (m *Middleware) Logger() {
	m.webServer.Use(fiberLogger.New(fiberLogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006/01/02 03:04 PM",
		Output:     os.Stdout,
	}))
}

func (m *Middleware) ContentSecurityPolicy() {
	csp := "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline'; " + // Hanya mengizinkan script dari server sendiri
		"style-src 'self' 'unsafe-inline' fonts.googleapis.com; " + // Mengizinkan CSS dari server sendiri, dan Google Fonts
		"font-src 'unsafe-inline' fonts.gstatic.com fonts.scalar.com; " + // Mengizinkan font dari server sendiri, Google Fonts, dan Scalar
		"img-src 'self' data:; " + // Mengizinkan gambar dari server sendiri atau inline data
		"frame-src 'none'; " + // Tidak mengizinkan frame
		"connect-src 'self'; " + // Mengizinkan koneksi hanya ke server sendiri
		"object-src 'none';" // Tidak mengizinkan pemuatan objek seperti Flash

	m.webServer.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentSecurityPolicy, csp)
		return c.Next()
	})
}
