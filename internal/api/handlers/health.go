package handlers

import "github.com/gofiber/fiber/v2"

type HealthHandler struct {
	appName    string
	appVersion string
}

func NewHealth(appName, appVersion string) *HealthHandler {
	return &HealthHandler{
		appName:    appName,
		appVersion: appVersion,
	}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}

func (h *HealthHandler) Root(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"app":     h.appName,
		"version": h.appVersion,
		"status":  "running",
	})
}
