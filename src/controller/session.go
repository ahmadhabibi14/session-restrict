package controller

import (
	"fmt"
	"net/http"
	"session-restrict/helper"
	"session-restrict/src/dto/request"
	"session-restrict/src/dto/response"
	"session-restrict/src/lib/logger"
	"session-restrict/src/repo/sessions"
	"session-restrict/src/service"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
)

type Session struct {
	srvSession *service.Session
}

func NewSession(app *fiber.App, srvSession *service.Session) {
	handler := &Session{srvSession}

	app.Route("/api/sessions", func(router fiber.Router) {
		router.Get("/", mustLoggedInAjax, handler.GetSessions)
		router.Patch("/approve", mustLoggedInAjax, handler.Approve)
		router.Patch("/delete", mustLoggedInAjax, handler.Delete)
	})
}

func (a *Session) Approve(c *fiber.Ctx) error {
	in, err := helper.ReadBody[request.ReqSessionApprove](c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ResponseCommon{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
	}

	session := getSession(c)

	out, err := a.srvSession.Approve(in, session.UserId)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	out.SetMessage(`Session Approved !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}

func (a *Session) Delete(c *fiber.Ctx) error {
	in, err := helper.ReadBody[request.ReqSessionDelete](c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.ResponseCommon{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
	}

	session := getSession(c)
	out, err := a.srvSession.Delete(in, session.UserId)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	out.SetMessage(`Session Deleted !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}

func (a *Session) GetSessions(c *fiber.Ctx) error {
	session := getSession(c)

	out, err := a.srvSession.GetSessions(session.UserId, session.Role)
	if err != nil {
		return c.Status(out.StatusCode).JSON(response.ResponseCommon{
			StatusCode: out.StatusCode,
			Error:      err.Error(),
		})
	}

	out.SetMessage(`Sessions obtained !`)
	out.SetStatus(http.StatusOK)

	return c.Status(http.StatusOK).JSON(out)
}

const (
	DeviceDesktop = "desktop"
	DeviceMobile  = "mobile"
	DeviceBot     = "bot"
	DeviceUnknown = "unknown"
)

const CookieAccessToken = `access_token`

func mustLoggedIn(c *fiber.Ctx) error {
	accessToken := c.Cookies(CookieAccessToken)
	if accessToken == `` {
		return c.Redirect("/signin", http.StatusTemporaryRedirect)
	}

	sess := sessions.NewSession()
	sess.AccessToken = accessToken

	session, err := sess.GetSessionByToken()
	if err != nil {
		RemoveAuthCookie(c)
		c.ClearCookie(CookieAccessToken)
		return c.Redirect("/", http.StatusTemporaryRedirect)
	}

	if !session.Approved {
		return c.Render("forbidden", fiber.Map{
			`Title`: fmt.Sprintf("%d - %s", fiber.StatusForbidden, `🚫 Access Denied`),
		}, "_layout")
	}

	return c.Next()
}

func mustLoggedInAjax(c *fiber.Ctx) error {
	accessToken := c.Cookies(CookieAccessToken)
	if accessToken == `` {
		return c.Status(http.StatusUnauthorized).JSON(response.ResponseCommon{
			StatusCode: http.StatusUnauthorized,
			Error:      "unauthorized",
		})
	}

	sess := sessions.NewSession()
	sess.AccessToken = accessToken

	session, err := sess.GetSessionByToken()
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(response.ResponseCommon{
			StatusCode: http.StatusUnauthorized,
			Error:      "unauthorized",
		})
	}

	if !session.Approved {
		return c.Status(http.StatusUnauthorized).JSON(response.ResponseCommon{
			StatusCode: http.StatusUnauthorized,
			Error:      "unauthorized",
		})
	}

	return c.Next()
}

func mustLoggedInAjaxUnapproved(c *fiber.Ctx) error {
	accessToken := c.Cookies(CookieAccessToken)
	if accessToken == `` {
		return c.Status(http.StatusUnauthorized).JSON(response.ResponseCommon{
			StatusCode: http.StatusUnauthorized,
			Error:      "unauthorized",
		})
	}

	sess := sessions.NewSession()
	sess.AccessToken = accessToken

	_, err := sess.GetSessionByToken()
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(response.ResponseCommon{
			StatusCode: http.StatusUnauthorized,
			Error:      "unauthorized",
		})
	}

	return c.Next()
}

func getSession(c *fiber.Ctx) sessions.Session {
	accessToken := c.Cookies(CookieAccessToken)

	sess := sessions.NewSession()
	sess.AccessToken = accessToken

	session, err := sess.GetSessionByToken()
	if err != nil {
		logger.Log.Error(err)
	}

	return session
}

func mustLoggedOut(c *fiber.Ctx) error {
	accessToken := c.Cookies(CookieAccessToken)
	if accessToken == `` {
		return c.Next()
	}

	return c.Redirect(`/`, http.StatusTemporaryRedirect)
}

func SetAuthCookie(c *fiber.Ctx, tokenString string, expiredAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     CookieAccessToken,
		Value:    tokenString,
		Expires:  expiredAt,
		SameSite: "Lax",
		Secure:   false,
		HTTPOnly: true,
		Path:     `/`,
	})
}

func RemoveAuthCookie(c *fiber.Ctx) {
	c.ClearCookie(CookieAccessToken)
	c.Cookie(&fiber.Cookie{
		Name:     CookieAccessToken,
		Value:    "",
		Expires:  time.Date(-1, 0, 0, 0, 0, 0, 0, time.Local),
		SameSite: "Lax",
		Secure:   false,
		HTTPOnly: true,
		Path:     `/`,
	})
}

func isBot(userAgent string) bool {
	// List of common bot identifiers
	botKeywords := []string{
		"bot", "crawler", "spider", "slurp", "fetch", "curl", "wget", "python", "java", "httpclient",
		"xhr", "facebook", "twitter", "linkedin", "google", "bing", "yandex", "baidu", "pinterest",
		"duckduckgo", "bingbot", "yandexbot", "googlebot", "facebookexternalhit", "twitterbot",
		"sogou", "slurp", "ccbot", "yeti", "ahrefsbot", "semrushbot", "rogerbot", "cognitiveseo",
	}

	userAgent = strings.ToLower(userAgent)

	for _, keyword := range botKeywords {
		if strings.Contains(userAgent, keyword) {
			return true
		}
	}

	return false
}

func GetIpV4(c *fiber.Ctx) string {
	ip := c.IP()
	ips := c.IPs()
	if len(ips) > 0 {
		if ips[0] != "" {
			ip = ips[0]
		}
	}

	return ip
}

func GetIpV6(c *fiber.Ctx) string {
	ip := c.IP()
	ips := c.IPs()
	if len(ips) >= 2 {
		if ips[1] != "" {
			ip = ips[1]
		}
	}

	return ip
}

func GetDevice(c *fiber.Ctx) string {
	ua := useragent.New(c.Get(fiber.HeaderUserAgent))

	device := DeviceDesktop

	if ua.Mobile() {
		device = DeviceMobile
	}

	if ua.Bot() || isBot(c.Get(fiber.HeaderUserAgent)) {
		device = DeviceBot
	}

	if !(ua.Bot() || isBot(c.Get(fiber.HeaderUserAgent))) && ua.OS() == `` {
		device = DeviceUnknown
	}

	return device
}

func GetOS(c *fiber.Ctx) string {
	ua := useragent.New(c.Get(fiber.HeaderUserAgent))
	return ua.OS()
}
