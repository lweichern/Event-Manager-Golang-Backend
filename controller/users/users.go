package users

import (
	"github.com/gofiber/fiber/v2"
)

func AllUsers(c *fiber.Ctx) error {

	return c.SendString("All Users")
}
