package bootstrap

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func InitFiber(logger zerolog.Logger) *fiber.App {
	r := fiber.New(fiber.Config{
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})

	p := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(
		InitPrometheusMetrics(),
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		}),
	)

	r.Get("/metrics", func(c *fiber.Ctx) error {
		p(c.Context())
		return nil
	})

	// gofiber recover => https://docs.gofiber.io/api/middleware/recover
	r.Use(recover.New())

	return r
}
