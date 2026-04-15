package routes

import (
	"job-portal/controllers"
	"job-portal/middleware"
	"job-portal/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup registers all application routes
func Setup(
	router *gin.Engine,
	authCtrl *controllers.AuthController,
	jobCtrl *controllers.JobController,
	appCtrl *controllers.ApplicationController,
	userCtrl *controllers.UserController,
) {
	// ─── Swagger UI ───────────────────────────────────────────────────────────
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "job-portal"})
	})

	api := router.Group("/api/v1")

	// ─── Auth (public) ────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}

	// ─── Protected routes ─────────────────────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.JWTAuth())

	// ─── Jobs (public read, role-restricted write) ────────────────────────────
	jobs := protected.Group("/jobs")
	{
		// Candidate & Recruiter & Admin can view jobs
		jobs.GET("", jobCtrl.GetAllJobs)
		jobs.GET("/:id", jobCtrl.GetJobByID)

		// Candidate: apply for a job
		jobs.POST("/:id/apply",
			middleware.RequireRole(models.RoleCandidate),
			appCtrl.ApplyForJob,
		)

		// Recruiter: manage own jobs
		jobs.POST("",
			middleware.RequireRole(models.RoleRecruiter),
			jobCtrl.CreateJob,
		)
		jobs.PUT("/:id",
			middleware.RequireRole(models.RoleRecruiter),
			jobCtrl.UpdateJob,
		)
		jobs.DELETE("/:id",
			middleware.RequireRole(models.RoleRecruiter, models.RoleAdmin),
			jobCtrl.DeleteJob,
		)
	}

	// ─── Applications ─────────────────────────────────────────────────────────
	applications := protected.Group("/applications")
	{
		// Candidate: view own applications
		applications.GET("",
			middleware.RequireRole(models.RoleCandidate),
			appCtrl.GetMyApplications,
		)

		// Recruiter: view applicants for a job & update status
		applications.GET("/job/:id",
			middleware.RequireRole(models.RoleRecruiter),
			appCtrl.GetApplicationsByJob,
		)
		applications.PATCH("/:id/status",
			middleware.RequireRole(models.RoleRecruiter),
			appCtrl.UpdateApplicationStatus,
		)
	}

	// ─── Admin ────────────────────────────────────────────────────────────────
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(models.RoleAdmin))
	{
		admin.GET("/users", userCtrl.GetAllUsers)
		admin.DELETE("/users/:id", userCtrl.DeleteUser)
		admin.GET("/jobs", jobCtrl.AdminGetAllJobs)
		admin.GET("/applications", appCtrl.AdminGetAllApplications)
	}
}
