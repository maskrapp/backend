package user

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maskrapp/backend/models"
	dbmodels "github.com/maskrapp/common/models"
	"gorm.io/gorm"
)

func DeleteMask(db *gorm.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		body := make(map[string]string)
		err := json.Unmarshal(c.Body(), &body)
		if err != nil {
			return c.SendStatus(500)
		}
		val, ok := body["mask"]
		if !ok {
			return c.Status(400).JSON(&models.APIResponse{
				Success: false,
				Message: "Invalid Body",
			})
		}
		mask := strings.ToLower(val)
		userID := c.Locals("user_id").(string)

		err = db.Delete(&dbmodels.Mask{}, "mask = ? AND user_id = ?", mask, userID).Error
		if err != nil {
			return c.Status(500).JSON(&models.APIResponse{
				Success: false,
				Message: "Something went wrong!",
			})
		}
		return c.JSON(&models.APIResponse{
			Success: true,
			Message: "Mask deleted",
		})
	}
}
