package fiber

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CORSConfig holds CORS middleware configuration for Fiber
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           3600,
	}
}

// CORS creates a CORS middleware for Fiber with the given configuration
func CORS(config CORSConfig) fiber.Handler {
	// Set defaults if not provided
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"}
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	}
	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	}
	if config.MaxAge == 0 {
		config.MaxAge = 3600
	}

	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")
	exposeHeaders := strings.Join(config.ExposeHeaders, ", ")

	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Check if origin is allowed
		if origin != "" {
			if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] == "*" {
				c.Set("Access-Control-Allow-Origin", "*")
			} else {
				for _, allowedOrigin := range config.AllowOrigins {
					if allowedOrigin == origin {
						c.Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		}

		// Set CORS headers
		c.Set("Access-Control-Allow-Methods", allowMethods)
		c.Set("Access-Control-Allow-Headers", allowHeaders)

		if len(config.ExposeHeaders) > 0 {
			c.Set("Access-Control-Expose-Headers", exposeHeaders)
		}

		if config.AllowCredentials {
			c.Set("Access-Control-Allow-Credentials", "true")
		}

		if config.MaxAge > 0 {
			c.Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
