package server

import (
	"time"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/httpapi"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func New(cfg *config.Config, store *store.Store) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "trip-listik-api",
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorHandler: jsonErrorHandler,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} ${latency} ${method} ${path} ${ip} ${locals:requestid}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowHeaders: "Authorization, Content-Type",
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
	}))

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/readyz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ready"})
	})

	v1 := app.Group("/v1")
	httpapi.Register(v1, cfg, store)

	return app
}

func jsonErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"message": err.Error(),
			"code":    code,
		},
	})
}
