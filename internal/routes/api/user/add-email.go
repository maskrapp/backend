package user

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/maskrapp/api/internal/global"
	"github.com/maskrapp/api/internal/models"
	"github.com/maskrapp/api/internal/utils"
)

func AddEmail(ctx global.Context) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := make(map[string]string)
		err := json.Unmarshal(c.Body(), &body)
		if err != nil {
			return c.SendStatus(500)
		}
		email, ok := body["email"]
		if !ok {
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
				Message: "Invalid body",
			})
		}
		if !utils.EmailRegex.MatchString(email) {
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
				Message: "Invalid email",
			})
		}

		domain := strings.Split(email, "@")[1]

		if _, err := ctx.Instances().Domains.Get(domain); err == nil {
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
				Message: "You cannot use that email",
			})
		}

		userId := c.Locals("user_id").(string)

		var result struct {
			Found bool
		}

		db := ctx.Instances().Gorm

		db.Raw("SELECT EXISTS(SELECT 1 FROM emails WHERE user_id = ? AND email = ?) AS found",
			userId, email).Scan(&result)

		if result.Found {
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
				Message: "That email is already registered to your account",
			})
		}

		emailRecord := &models.Email{
			UserID:     userId,
			Email:      email,
			IsVerified: false,
			IsPrimary:  false,
		}

		err = db.Create(emailRecord).Error

		if err != nil {
			return c.Status(500).JSON(&models.APIResponse{
				Success: false,
				Message: "Something went wrong!",
			})
		}
		emailVerification := &models.EmailVerification{
			EmailID:          emailRecord.Id,
			VerificationCode: utils.GenerateCode(5),
			ExpiresAt:        time.Now().Add(30 * time.Minute).Unix(),
		}
		err = db.Create(emailVerification).Error

		if err != nil {
			return c.Status(500).JSON(&models.APIResponse{
				Success: false,
				Message: "Something went wrong!",
			})
		}
		return c.JSON(&models.APIResponse{
			Success: true,
		})
	}
}
