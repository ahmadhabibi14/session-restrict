package controller

import "github.com/gofiber/fiber/v2"

type Pages struct{}

func NewPages(app *fiber.App) {
	handler := &Pages{}

	app.Get("/", mustLoggedIn, handler.Home)
	app.Get("/signin", mustLoggedOut, handler.SignIn)
	app.Get("/signup", mustLoggedOut, handler.SignUp)
}

func (p *Pages) Home(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	c.Set(fiber.HeaderXFrameOptions, "SAMEORIGIN")

	return c.Render("index", fiber.Map{
		`Title`: "Session Restriction",
	}, "_layout")
}

func (p *Pages) SignIn(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	c.Set(fiber.HeaderXFrameOptions, "SAMEORIGIN")

	return c.Render("signin", fiber.Map{
		`Title`: "Session Restriction",
	}, "_layout")
}

func (p *Pages) SignUp(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	c.Set(fiber.HeaderXFrameOptions, "SAMEORIGIN")

	return c.Render("signup", fiber.Map{
		`Title`: "Session Restriction",
	}, "_layout")
}
