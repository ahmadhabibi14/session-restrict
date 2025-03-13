package controller

import (
	"net/http"
	"session-restrict/src/repo/sessions"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
)

const (
	DeviceDesktop = "desktop"
	DeviceMobile  = "mobile"
	DeviceBot     = "bot"
	DeviceUnknown = "unknown"
)

func mustLoggedIn(c *fiber.Ctx) error {
	accessToken := c.Cookies(`auth`)
	if accessToken == `` {
		return c.Redirect("/signin", http.StatusPermanentRedirect)
	}

	_, err := sessions.GetSessionByToken(accessToken)
	if err != nil {
		RemoveAuthCookie(c)
		return c.Redirect("/", http.StatusPermanentRedirect)
	}

	return c.Next()
}

func mustLoggedOut(c *fiber.Ctx) error {
	accessToken := c.Cookies(`auth`)
	if accessToken == `` {
		RemoveAuthCookie(c)
		return c.Next()
	}

	return c.Redirect(`/`, http.StatusPermanentRedirect)
}

func SetAuthCookie(c *fiber.Ctx, tokenString string, expiredAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     `auth`,
		Value:    tokenString,
		Expires:  expiredAt,
		SameSite: "Lax",
		Secure:   false,
		HTTPOnly: true,
		Path:     `/`,
	})
}

func RemoveAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:    `auth`,
		Value:   "",
		Path:    `/`,
		Expires: time.Date(-1, 0, 0, 0, 0, 0, 0, time.Local),
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
