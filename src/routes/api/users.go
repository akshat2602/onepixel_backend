package api

import (
	"errors"
	"onepixel_backend/src/controllers"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/security"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

var usersController *controllers.UsersController

// UsersRoute defines the routes for /api/v1/users
func UsersRoute(db *gorm.DB) func(router fiber.Router) {
	usersController = controllers.NewUsersController(db)
	return func(router fiber.Router) {
		router.Post("/", registerUser)
		router.Post("/login", loginUser)
		router.Get("/:id", security.MandatoryAuthMiddleware, getUserInfo)
		router.Patch("/:id", security.MandatoryAuthMiddleware, updateUserInfo)
	}
}

// registerUser
//
//	@Summary		Register new user
//	@Description	Register new user
//	@ID				register-user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dtos.CreateUserRequest	true	"User"
//	@Success		201		{object}	dtos.UserResponse
//	@Failure		400		{object}	dtos.ErrorResponse "The request body is not valid"
//	@Failure		422		{object}	dtos.ErrorResponse "email and password are required to create user"
//	@Failure		409		{object}	dtos.ErrorResponse "User with this email already exists"
//	@Router			/api/v1/users [post]
//	@Security		ApiKeyAuth
func registerUser(ctx *fiber.Ctx) error {
	var u = new(dtos.CreateUserRequest)
	if err := ctx.BodyParser(u); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.CreateErrorResponse(
			fiber.StatusBadRequest,
			"The request body is not valid",
		))
	}

	if u.Email == "" || u.Password == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(dtos.CreateErrorResponse(
			fiber.StatusUnprocessableEntity,
			"email and password are required to create user",
		))
	}

	savedUser, token, err := usersController.Create(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ctx.Status(fiber.StatusConflict).JSON(dtos.CreateErrorResponse(fiber.StatusConflict, "User with this email already exists"))
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.CreateUserResponseFromUser(savedUser, &token))
}

// loginUser
//
// @Summary		Login user
// @Description	Login user
// @ID				login-user
// @Tags			users
// @Accept			json
// @Produce		json
// @Router			/api/v1/users/login [post]
// @Security		ApiKeyAuth
func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("LoginUser")
}

// getUserInfo
//
// @Summary		Get user info
// @Description	Get user info
// @ID				get-user-info
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path	uint	true	"User ID"
// @Router			/api/v1/users/:id [get]
func getUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("GetUserInfo")
}

// updateUserInfo
//
// @Summary		Update user info
// @Description	Update user info
// @ID				update-user-info
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path	uint	true	"User ID"
// @Router			/api/v1/users/:id [patch]
func updateUserInfo(ctx *fiber.Ctx) error {
	return ctx.SendString("UpdateUserInfo")
}
