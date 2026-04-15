package controllers

import (
	"job-portal/services"
	"job-portal/utils"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new account with name, email, password and role
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.RegisterInput  true  "Register payload"
// @Success      201   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      409   {object}  utils.APIResponse
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	result, err := ctrl.authService.Register(input)
	if err != nil {
		if err.Error() == "email already registered" {
			utils.Conflict(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "user registered successfully", result)
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password, returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.LoginInput  true  "Login payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	result, err := ctrl.authService.Login(input)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, "login successful", result)
}
