package controller

import (
	"github.com/gofiber/fiber/v2"
	"infinitoon.dev/infinitoon/apps/web/service"
	appctx "infinitoon.dev/infinitoon/pkg/context"
	"infinitoon.dev/infinitoon/pkg/rest"
)

type UserController interface {
	rest.Controller
	GetByID(ctx *fiber.Ctx) (err error)
	GetByEmail(ctx *fiber.Ctx) (err error)
}

type userController struct {
	appCtx      *appctx.AppContext
	userService service.UserService
}

func NewUserController(appCtx *appctx.AppContext, userService service.UserService) UserController {
	return &userController{
		appCtx:      appCtx,
		userService: userService,
	}
}

func (c *userController) Route() *rest.RestRoute {
	return rest.NewRestRoute().
		SetPrefix("/users").
		Handler(func(router fiber.Router) {
			router.Get("/:id", c.GetByID)
			router.Get("/email/:email", c.GetByEmail)
		})
}

func (c *userController) GetByID(ctx *fiber.Ctx) (err error) {
	id := ctx.Params("id")
	user, err := c.userService.GetByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(user)
}

func (c *userController) GetByEmail(ctx *fiber.Ctx) (err error) {
	email := ctx.Params("email")
	err = c.appCtx.GetValidator().Var(email, "required,email")
	if err != nil {
		return err
	}
	user, err := c.userService.GetByEmail(email)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.JSON(user)
}
