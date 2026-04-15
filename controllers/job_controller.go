package controllers

import (
	"job-portal/middleware"
	"job-portal/repository"
	"job-portal/services"
	"job-portal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JobController handles job-related endpoints
type JobController struct {
	jobService services.JobService
}

// NewJobController creates a new JobController
func NewJobController(jobService services.JobService) *JobController {
	return &JobController{jobService: jobService}
}

// GetAllJobs godoc
// @Summary      List all jobs
// @Description  Get paginated list of jobs with optional location and company filters
// @Tags         Jobs
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number"       default(1)
// @Param        page_size  query     int     false  "Items per page"    default(10)
// @Param        location   query     string  false  "Filter by location"
// @Param        company    query     string  false  "Filter by company"
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /jobs [get]
func (ctrl *JobController) GetAllJobs(c *gin.Context) {
	pagination := utils.GetPagination(c)
	filter := repository.JobFilter{
		Location: c.Query("location"),
		Company:  c.Query("company"),
	}

	result, err := ctrl.jobService.GetAllJobs(pagination, filter)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "jobs retrieved successfully", result)
}

// GetJobByID godoc
// @Summary      Get a job by ID
// @Description  Retrieve a single job by its UUID
// @Tags         Jobs
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Job UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /jobs/{id} [get]
func (ctrl *JobController) GetJobByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid job ID", nil)
		return
	}

	job, err := ctrl.jobService.GetJobByID(id)
	if err != nil {
		if err.Error() == "job not found" {
			utils.NotFound(c, err.Error())
			return
		}
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "job retrieved successfully", job)
}

// CreateJob godoc
// @Summary      Create a job
// @Description  Recruiter creates a new job posting
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      services.CreateJobInput  true  "Job payload"
// @Success      201   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Failure      403   {object}  utils.APIResponse
// @Router       /jobs [post]
func (ctrl *JobController) CreateJob(c *gin.Context) {
	var input services.CreateJobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	recruiterID := middleware.GetUserID(c)
	job, err := ctrl.jobService.CreateJob(input, recruiterID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, "job created successfully", job)
}

// UpdateJob godoc
// @Summary      Update a job
// @Description  Recruiter updates their own job posting
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                   true  "Job UUID"
// @Param        body  body      services.UpdateJobInput  true  "Update payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      403   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /jobs/{id} [put]
func (ctrl *JobController) UpdateJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid job ID", nil)
		return
	}

	var input services.UpdateJobInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, "validation error", err.Error())
		return
	}

	recruiterID := middleware.GetUserID(c)
	job, err := ctrl.jobService.UpdateJob(id, input, recruiterID)
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

	utils.Success(c, "job updated successfully", job)
}

// DeleteJob godoc
// @Summary      Delete a job
// @Description  Recruiter deletes their own job (Admin can delete any)
// @Tags         Jobs
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Job UUID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /jobs/{id} [delete]
func (ctrl *JobController) DeleteJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid job ID", nil)
		return
	}

	recruiterID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	if err := ctrl.jobService.DeleteJob(id, recruiterID, role); err != nil {
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

	utils.Success(c, "job deleted successfully", nil)
}

// AdminGetAllJobs godoc
// @Summary      [Admin] List all jobs
// @Description  Admin retrieves all job postings with pagination and filters
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int     false  "Page number"     default(1)
// @Param        page_size  query     int     false  "Items per page"  default(10)
// @Param        location   query     string  false  "Filter by location"
// @Param        company    query     string  false  "Filter by company"
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Router       /admin/jobs [get]
func (ctrl *JobController) AdminGetAllJobs(c *gin.Context) {
	pagination := utils.GetPagination(c)
	filter := repository.JobFilter{
		Location: c.Query("location"),
		Company:  c.Query("company"),
	}

	result, err := ctrl.jobService.GetAllJobs(pagination, filter)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "jobs retrieved successfully", result)
}
