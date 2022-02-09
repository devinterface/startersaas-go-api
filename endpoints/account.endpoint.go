package endpoints

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

// AccountEndpoint struct
type AccountEndpoint struct{ BaseEndpoint }

// ByID function
func (accountEndpoint *AccountEndpoint) ByID(ctx *fiber.Ctx) error {
	account, err := accountService.ByID(ctx.Params("id"))
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showAccount := models.ShowAccountSerializer().Transform(account)
	return ctx.JSON(showAccount)
}

// Update function
func (accountEndpoint *AccountEndpoint) Update(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"companyName":           "ascii",
		"companyVat":            "ascii",
		"companyBillingAddress": "ascii",
		"companySdi":            "ascii",
		"companyPhone":          "ascii",
		"companyEmail":          "ascii",
		"companyPec":            "ascii",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	updatedAccount, err := accountService.Update(ctx.Params("id"), inputMap)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showAccount := models.ShowAccountSerializer().Transform(updatedAccount)
	return ctx.JSON(showAccount)
}
