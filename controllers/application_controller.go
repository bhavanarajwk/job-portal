package controllers

import (
	"job-portal/middleware"
	"job-portal/services"
	"job-portal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationController handles application-related endpoints
type ApplicationController struct {
	appService services.ApplicationService
}

// NewApplicationController creates a new ApplicationController
func NewApplicationController(appService services.ApplicationService) *ApplicationController {
	return &ApplicationController{appService: appService}
}

// ApplyForJob godoc
// @Summary      Apply for a job
// @Description  Candidate submits an application for a job
// @Tags         Applications
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Job UUID"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Router       /jobs/{id}/apply [post]
func (ctrl *ApplicationController) ApplyForJob(c *gin.Context) {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid job ID", nil)
		return
	}

	candidateID := middleware.GetUserID(c)
	app, err := ctrl.appService.Apply(candidateID, jobID)
	if err != nil {
		switch err.Error() {
		case "job not found":
			utils.NotFound(c, err.Error())
		case "you have already applied for this job":
			utils.Conflict(c, err.Error())
		default:
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	utils.Created(c, "application submitted successfully", app)
}

// GetMyApplications godoc
// @Summary      My applications
// @Description  Candidate views their own job applications
// @Tags         Applications
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Page number"     default(1)
// @Param        page_size  query     int  false  "Items per page"  default(10)
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Router       /applications [get]
func (ctrl *ApplicationController) GetMyApplications(c *gin.Context) {
	candidateID := middleware.GetUserID(c)
	pagination := utils.GetPagination(c)

	result, err := ctrl.appService.GetMyApplications(candidateID, pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "applications retrieved successfully", result)
}

// GetApplicationsByJob godoc
// @Summary      Applicants for a job
// @Description  Recruiter views all applicants for one of their jobs
// @Tags         Applications
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string  true   "Job UUID"
// @Param        page       query     int     false  "Page number"     default(1)
// @Param        page_size  query     int     false  "Items per page"  default(10)
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /applications/job/{id} [get]
func (ctrl *ApplicationController) GetApplicationsByJob(c *gin.Context) {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid job ID", nil)
		return
	}

	recruiterID := middleware.GetUserID(c)
	pagination := utils.GetPagination(c)

	result, err := ctrl.appService.GetApplicationsByJob(jobID, recruiterID, pagination)
	if err != nil {
		switch err.Error() {
		case "job not found":
			utils.NotFound(c, err.Error())
		case "forbidden: you do not own this job":
			utils.Forbidden(c, err.Error())
		default:
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	utils.Success(c, "applicants retrieved successfully", result)
}

// UpdateApplicationStatus godoc
// @Summary      Update application status
// @Description  Recruiter updates the status of an application (REVIEWED, ACCEPTED, REJECTED)
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                      true  "Application UUID"
// @Param        body  body      services.UpdateStatusInput  true  "Status payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      403   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /applications/{id}/status [patch]
func (ctrl *ApplicationController) UpdateApplicationStatus(c *gin.Context) {
	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid application ID", nil)
		return
	}

	var input services.UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	recruiterID := middleware.GetUserID(c)
	app, err := ctrl.appService.UpdateStatus(appID, recruiterID, input)
	if err != nil {
		switch err.Error() {
		case "application not found":
			utils.NotFound(c, err.Error())
		case "forbidden: you do not own this job":
			utils.Forbidden(c, err.Error())
		default:
			utils.InternalServerError(c, err.Error())
		}
		return
	}

	utils.Success(c, "application status updated", app)
}

// AdminGetAllApplications godoc
// @Summary      [Admin] List all applications
// @Description  Admin retrieves all applications across all jobs
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "Page number"     default(1)
// @Param        page_size  query     int  false  "Items per page"  default(10)
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Router       /admin/applications [get]
func (ctrl *ApplicationController) AdminGetAllApplications(c *gin.Context) {
	pagination := utils.GetPagination(c)

	result, err := ctrl.appService.GetAllApplications(pagination)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "applications retrieved successfully", result)
}
