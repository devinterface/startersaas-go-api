package endpoints

import (
	"devinterface.com/startersaas-go-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"go.mongodb.org/mongo-driver/bson"
)

// TeamEndpoint struct
type TeamEndpoint struct{ BaseEndpoint }

func (teamEndpoint *TeamEndpoint) Index(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	params := bson.M{}
	ctx.BodyParser(&params)
	params["accountId"] = me.AccountID
	teams, err := teamService.FindBy(params)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showTeams, _ := models.ShowTeamSerializer().TransformArray(teams)
	return ctx.JSON(showTeams)
}

func (teamEndpoint *TeamEndpoint) ByID(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	team, err := teamService.ByID(ctx.Params("id"), me.AccountID)
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showTeam := models.ShowTeamSerializer().Transform(team)
	return ctx.JSON(showTeam)
}

// Create function
func (teamEndpoint *TeamEndpoint) Create(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	currentAccount, _ := userEndpoint.CurrentAccount(ctx)

	teams, err := teamService.FindBy(bson.M{"accountId": me.AccountID})
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if len(teams) >= models.MaxTeamsPerPlan(currentAccount.PlanType) {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You have reached the maximum number of teams for your plan. Upgrade it to increase the number.",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	v := validate.Map(inputMap)
	v.StringRule("code", "string|required")
	v.StringRule("name", "string|required")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}
	team, err := teamService.Create(v.SafeData(), me.AccountID)
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showTeam := models.ShowTeamSerializer().Transform(team)
	return ctx.JSON(showTeam)
}

// Update function
func (teamEndpoint *TeamEndpoint) Update(ctx *fiber.Ctx) error {
	me, _ := teamEndpoint.CurrentUser(ctx)
	if can := teamEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}

	var inputMap = make(map[string]interface{})
	ctx.BodyParser(&inputMap)

	v := validate.Map(inputMap)
	v.StringRule("name", "string")

	if !v.Validate() {
		return ctx.Status(422).JSON(v.Errors)
	}

	updateTeam, err := teamService.Update(ctx.Params("id"), me.AccountID, v.SafeData())
	if err != nil {
		return ctx.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	showTeam := models.ShowTeamSerializer().Transform(updateTeam)
	return ctx.JSON(showTeam)
}

// Delete function
func (teamEndpoint *TeamEndpoint) Delete(ctx *fiber.Ctx) error {
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	deleted, err := teamService.Delete(ctx.Params("id"))
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"deleted": deleted,
	})
}

// Add user function
func (teamEndpoint *TeamEndpoint) AddUser(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	team, err := teamService.AddUser(ctx.Params("id"), me.AccountID, ctx.Params("userId"))
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showTeam := models.ShowTeamSerializer().Transform(team)
	return ctx.JSON(showTeam)
}

// REmove user function
func (teamEndpoint *TeamEndpoint) RemoveUser(ctx *fiber.Ctx) error {
	me, _ := userEndpoint.CurrentUser(ctx)
	if can := userEndpoint.Can(ctx, models.AdminRole); !can {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "You are not authorized to perform this action",
		})
	}
	team, err := teamService.RemoveUser(ctx.Params("id"), me.AccountID, ctx.Params("userId"))
	if err != nil {
		return ctx.Status(422).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	showTeam := models.ShowTeamSerializer().Transform(team)
	return ctx.JSON(showTeam)
}
