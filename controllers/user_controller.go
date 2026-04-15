package controllers

import (
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
