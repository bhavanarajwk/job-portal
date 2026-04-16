// @title           Job Portal API
// @version         1.0
// @description     A production-ready Job Portal backend with role-based access control.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@jobportal.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token as: Bearer <token>

package main

import (
	"fmt"
	"log"

	"job-portal/config"
	"job-portal/controllers"
	"job-portal/database"
	docs "job-portal/docs"
	"job-portal/middleware"
	"job-portal/repository"
	"job-portal/routes"
	"job-portal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// ── 1. Load configuration ─────────────────────────────────────────────────
	cfg := config.Load()

	// ── 2. Connect to database & run migrations ───────────────────────────────
	database.Connect(cfg)

	// ── 3. Wire up repositories ───────────────────────────────────────────────
	userRepo := repository.NewUserRepository(database.DB)
	jobRepo := repository.NewJobRepository(database.DB)
	appRepo := repository.NewApplicationRepository(database.DB)

	// ── 4. Wire up services ───────────────────────────────────────────────────
	authService := services.NewAuthService(userRepo)
	jobService := services.NewJobService(jobRepo, userRepo, appRepo)
	appService := services.NewApplicationService(appRepo, jobRepo, userRepo)
	userService := services.NewUserService(userRepo)

	// ── 5. Wire up controllers ────────────────────────────────────────────────
	authCtrl := controllers.NewAuthController(authService)
	jobCtrl := controllers.NewJobController(jobService)
	appCtrl := controllers.NewApplicationController(appService)
	userCtrl := controllers.NewUserController(userService)

	// ── 6. Set up Gin router ──────────────────────────────────────────────────
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// ── 7. Update swagger host to match configured port ──────────────────────
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.AppPort)

	// ── 8. Register routes ────────────────────────────────────────────────────
	routes.Setup(router, authCtrl, jobCtrl, appCtrl, userCtrl)

	// ── 8. Start server ───────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[INFO] Job Portal API starting on %s (env: %s)", addr, cfg.AppEnv)
	log.Printf("[INFO] Swagger UI available at http://localhost%s/swagger/index.html", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[FATAL] Server failed to start: %v", err)
	}
}
