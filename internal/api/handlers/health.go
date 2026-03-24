package handlers

import "github.com/gofiber/fiber/v2"

// HealthHandler serves health-check and informational endpoints.
type HealthHandler struct {
	appName    string
	appVersion string
}

// NewHealth creates a new HealthHandler.
func NewHealth(appName, appVersion string) *HealthHandler {
	return &HealthHandler{
		appName:    appName,
		appVersion: appVersion,
	}
}

// Health returns a simple health-check response.
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// Live returns a liveness probe response.
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

// Ready returns a readiness probe response.
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}

// Root returns application metadata.
func (h *HealthHandler) Root(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"app":     h.appName,
		"version": h.appVersion,
		"status":  "running",
	})
}
