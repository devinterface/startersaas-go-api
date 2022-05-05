package endpoints

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
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

	v := validate.Map(inputMap)
	v.StringRule("companyName", "ascii")
	v.StringRule("companyVat", "ascii")
	v.StringRule("companyBillingAddress", "ascii")
	v.StringRule("companySdi", "ascii")
	v.StringRule("companyEmail", "email")
	v.StringRule("companyPec", "email")
	v.StringRule("companyCountry", "ascii")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
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
