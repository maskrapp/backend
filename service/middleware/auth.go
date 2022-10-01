package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maskrapp/backend/jwt"
	"github.com/maskrapp/backend/models"
)

func AuthMiddleware(jwtHandler *jwt.JWTHandler) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		auth := c.GetReqHeaders()["Authorization"]
		if auth == "" {
			return c.SendStatus(401)
		}
		split := strings.Split(auth, " ")
		accessToken := split[1]
		claims, err := jwtHandler.Validate(accessToken, false)
		if err != nil {
			fmt.Println("err", err)
			if strings.Contains(err.Error(), "token mismatch") {
				return c.Status(400).JSON(&models.APIResponse{
					Success: false,
					Message: "Token mismatch",
				})
			}
			if strings.Contains(err.Error(), "token is expired by") {
				return c.Status(401).JSON(&models.APIResponse{
					Success: false,
					Message: "Token is expired",
				})
			}
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
			})

		}
		c.Locals("user_id", claims.UserId)
		return c.Next()
	}
}