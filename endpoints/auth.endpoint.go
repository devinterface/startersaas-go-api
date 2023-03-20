package endpoints

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
)

// AuthEndpoint struct
type AuthEndpoint struct{ BaseEndpoint }

// Login function
func (authEndpoint *AuthEndpoint) Login(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("email", "email|required")
	v.StringRule("password", "required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := authService.Login(inputMap["email"].(string), inputMap["password"].(string), false)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})

	}
	return ctx.JSON(response)
}

// SsoLogin function
func (authEndpoint *AuthEndpoint) SsoLogin(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("sso", "required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := authService.Sso(inputMap["sso"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})

	}
	return ctx.JSON(response)
}

// Signup function
func (authEndpoint *AuthEndpoint) Signup(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("subdomain", "required")
	v.StringRule("email", "email|required")
	v.StringRule("password", "required")
	v.StringRule("privacyAccepted", "required")
	v.StringRule("marketingAccepted", "-")
	v.StringRule("language", "-")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}

	response, err := authService.Signup(inputMap, os.Getenv("SIGNUP_WITH_ACTIVATE") == "true")
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// SendActivationLink function
func (authEndpoint *AuthEndpoint) SendActivationLink(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("email", "email|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := emailService.SendActivationEmail(bson.M{"email": inputMap["email"].(string)})
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// Activate function
func (authEndpoint *AuthEndpoint) Activate(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("email", "email|required")
	v.StringRule("token", "required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := authService.Activate(inputMap["token"].(string), inputMap["email"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// SendForgotPasswordLink function
func (authEndpoint *AuthEndpoint) SendForgotPasswordLink(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("email", "email|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := emailService.SendForgotPasswordEmail(bson.M{"email": inputMap["email"].(string)})
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// ResetPassword function
func (authEndpoint *AuthEndpoint) ResetPassword(ctx *fiber.Ctx) error {
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	v := validate.Map(inputMap)
	v.StringRule("email", "email|required")
	v.StringRule("password", "ascii|required")
	v.StringRule("passwordResetToken", "ascii|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	response, err := authService.ResetPassword(inputMap["passwordResetToken"].(string), inputMap["password"].(string), inputMap["email"].(string))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}

// RefreshToken function
func (authEndpoint *AuthEndpoint) RefreshToken(ctx *fiber.Ctx) error {
	me, err := userEndpoint.CurrentUser(ctx)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	response, err := authService.Login(me.Email, "", true)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(response)
}
