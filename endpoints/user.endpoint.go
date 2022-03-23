package endpoints

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// UserEndpoint struct
type UserEndpoint struct{ BaseEndpoint }

// Me function
func (userEndpoint *UserEndpoint) Me(ctx *fiber.Ctx) error {
	me, err := userEndpoint.CurrentUser(ctx)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showUser := models.ShowUserSerializer().Transform(me)

	queryAccount := ctx.Query("withAccount")
	if queryAccount == "true" {
		account, err := accountService.ByID(me.AccountID)
		if err != nil {
			return ctx.Status(404).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		showAccount := models.ShowAccountSerializer().Transform(account)
		showUser["account"] = showAccount
	}

	return ctx.JSON(showUser)
}

// UpdateMe function
func (userEndpoint *UserEndpoint) UpdateMe(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"name":     "alpha",
		"surname":  "alpha",
		"language": "alpha",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	updatedUser, _ := userService.Update(me.GetID(), me.AccountID, inputMap)
	showUser := models.ShowUserSerializer().Transform(updatedUser)
	return ctx.JSON(showUser)
}

// ChangePassword function
func (userEndpoint *UserEndpoint) ChangePassword(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)
	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"password": "ascii,required",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	updatedUser, _ := userService.UpdatePassword(me.GetID(), inputMap["password"].(string))
	showUser := models.ShowUserSerializer().Transform(updatedUser)
	return ctx.JSON(showUser)
}

// GenerateSso function
func (userEndpoint *UserEndpoint) GenerateSso(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	ssoUUID, _ := uuid.NewRandom()
	inputMap := map[string]string{"sso": ssoUUID.String()}

	updatedUser, err := userService.Update(me.GetID(), me.AccountID, inputMap)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"sso": updatedUser.Sso,
	})
}

// Index function
func (userEndpoint *UserEndpoint) Index(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	params := bson.M{}
	ctx.BodyParser(&params)
	params["accountId"] = me.AccountID
	users, err := userService.FindBy(params)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showUsers, _ := models.ShowUserSerializer().TransformArray(users)
	return ctx.JSON(showUsers)
}

// Create function
func (userEndpoint *UserEndpoint) Create(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"name":     "alpha",
		"surname":  "alpha",
		"language": "alpha, in(it|en)",
		"email":    "email, required",
		"password": "alpha",
		"role":     "alpha, in(user|admin)",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	user, err := userService.Create(inputMap, me.AccountID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showUser := models.ShowUserSerializer().Transform(user)
	return ctx.JSON(showUser)
}

// Update function
func (userEndpoint *UserEndpoint) Update(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	_, err := govalidator.ValidateMap(inputMap, map[string]interface{}{
		"name":     "alpha",
		"surname":  "alpha",
		"language": "alpha, in(it|en)",
		"role":     "alpha, in(user|admin)",
	})
	if err != nil {
		return ctx.Status(422).JSON(err.Error())
	}
	updatedUser, err := userService.Update(ctx.Params("id"), me.AccountID, inputMap)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showUser := models.ShowUserSerializer().Transform(updatedUser)
	return ctx.JSON(showUser)
}

// Delete function
func (userEndpoint *UserEndpoint) Delete(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	can := userEndpoint.Can(ctx, models.AdminRole)
	can = can && me.ID.Hex() != ctx.Params("id")
	if !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	deleted, err := userService.Delete(ctx.Params("id"), me.AccountID)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"deleted": deleted,
	})
}
