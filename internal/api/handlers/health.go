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

// Health godoc
// @Summary      Liveness check
// @Description  Returns 200 with status:ok when the process is running.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string "status:ok"
// @Router       /health [get]
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// Live godoc
// @Summary      Liveness probe (Kubernetes-style)
// @Description  Returns 200 with status:alive when the process is running.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string "status:alive"
// @Router       /live [get]
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}

// Ready godoc
// @Summary      Readiness probe
// @Description  Returns 200 when the service is ready to accept traffic.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string "status:ready"
// @Router       /ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}

// Root godoc
// @Summary      Service info
// @Description  Returns app name, version, and running status.
// @Tags         Health
// @Produce      json
// @Success      200 {object} map[string]string "app, version, status"
// @Router       / [get]
func (h *HealthHandler) Root(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"app":     h.appName,
		"version": h.appVersion,
		"status":  "running",
	})
}
