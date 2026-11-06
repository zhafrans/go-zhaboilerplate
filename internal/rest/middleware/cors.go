package middleware

import "github.com/gofiber/fiber/v2"

// CORS will handle the CORS middleware
func CORS(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	if c.Method() == fiber.MethodOptions {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Next()
}
