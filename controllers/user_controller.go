package controllers

import (
	"job-portal/middleware"
	"job-portal/services"
	"job-portal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController handles admin user management endpoints
type UserController struct {
	userService services.UserService
}

// NewUserController creates a new UserController
func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetAllUsers godoc
// @Summary      [Admin] List all users
// @Description  Admin retrieves all registered users with pagination
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Page number"     default(1)
// @Param        page_size  query     int  false  "Items per page"  default(10)
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Router       /admin/users [get]
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.GetPagination(c)

	result, err := ctrl.userService.GetAllUsers(pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "users retrieved successfully", result)
}

// DeleteUser godoc
// @Summary      [Admin] Delete a user
// @Description  Admin soft-deletes a user by UUID
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /admin/users/{id} [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid user ID", nil)
		return
	}

	if err := ctrl.userService.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "user deleted successfully", nil)
}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Returns profile details of the authenticated user
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /users/me [get]
func (ctrl *UserController) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := ctrl.userService.GetMe(userID)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "profile retrieved successfully", user)
}

// UpdateMe godoc
// @Summary      Update current user profile
// @Description  Updates authenticated user's name, email and/or password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      services.UpdateMeInput  true  "Update profile payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Failure      409   {object}  utils.APIResponse
// @Router       /users/me [put]
// @Router       /users/me [patch]
func (ctrl *UserController) UpdateMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var input services.UpdateMeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	user, err := ctrl.userService.UpdateMe(userID, input)
	if err != nil {
		switch err.Error() {
		case "at least one field is required":
			utils.BadRequest(c, err.Error(), nil)
		case "user not found":
			utils.NotFound(c, err.Error())
		case "email already registered":
			utils.Conflict(c, err.Error())
		default:
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	utils.Success(c, "profile updated successfully", user)
}

// ChangePassword godoc
// @Summary      Change current user password
// @Description  Changes authenticated user's password by validating current password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      services.ChangePasswordInput  true  "Change password payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /users/me/password [put]
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var input services.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	err := ctrl.userService.ChangePassword(userID, input)
	if err != nil {
		switch err.Error() {
		case "invalid current password":
			utils.Unauthorized(c, err.Error())
		case "new password must be different from current password":
			utils.BadRequest(c, err.Error(), nil)
		case "user not found":
			utils.NotFound(c, err.Error())
		default:
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	utils.Success(c, "password changed successfully", nil)
}
