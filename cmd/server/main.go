package main

import (
	"log"

	"dormi-api/internal/config"
	"dormi-api/internal/database"
	"dormi-api/internal/dto"
	"dormi-api/internal/handler"
	"dormi-api/internal/middleware"
	"dormi-api/internal/model"
	"dormi-api/internal/repository"
	"dormi-api/internal/service"

	_ "dormi-api/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Dormi API
// @version 1.0
// @description 기숙사 관리 시스템 API - 학생 상벌점 관리 및 당직 일정 운영

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT 토큰 (Bearer {token})

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	pointRepo := repository.NewPointRepository(db)
	pointReasonRepo := repository.NewPointReasonRepository(db)
	dutyRepo := repository.NewDutyRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	seedAdmin(cfg, authService)
	studentService := service.NewStudentService(studentRepo)
	pointReasonService := service.NewPointReasonService(pointReasonRepo)
	pointService := service.NewPointService(pointRepo, studentRepo, pointReasonRepo)
	dutyService := service.NewDutyService(dutyRepo)
	auditService := service.NewAuditService(auditRepo)

	authHandler := handler.NewAuthHandler(authService, auditService)
	studentHandler := handler.NewStudentHandler(studentService, auditService)
	pointReasonHandler := handler.NewPointReasonHandler(pointReasonService, auditService)
	pointHandler := handler.NewPointHandler(pointService, auditService)
	dutyHandler := handler.NewDutyHandler(dutyService, auditService)
	auditHandler := handler.NewAuditHandler(auditService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/auth/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authService))
	{
		api.PATCH("/auth/password", authHandler.ChangePassword)

		users := api.Group("/users")
		users.Use(middleware.RequireAdmin())
		{
			users.GET("", authHandler.GetAllUsers)
			users.GET("/:id", authHandler.GetUserByID)
			users.POST("", authHandler.CreateUser)
			users.PUT("/:id", authHandler.UpdateUser)
			users.DELETE("/:id", authHandler.DeleteUser)
		}

		students := api.Group("/students")
		{
			students.GET("", studentHandler.GetAll)
			students.GET("/:id", studentHandler.GetByID)
			students.POST("", middleware.RequireAdminOrSupervisor(), studentHandler.Create)
			students.PUT("/:id", middleware.RequireAdminOrSupervisor(), studentHandler.Update)
			students.DELETE("/:id", middleware.RequireAdminOrSupervisor(), studentHandler.Delete)
			students.POST("/import", middleware.RequireAdminOrSupervisor(), studentHandler.Import)
		}

		pointReasons := api.Group("/point-reasons")
		{
			pointReasons.GET("", pointReasonHandler.GetAll)
			pointReasons.GET("/:id", pointReasonHandler.GetByID)
			pointReasons.POST("", middleware.RequireAdminOrSupervisor(), pointReasonHandler.Create)
			pointReasons.PUT("/:id", middleware.RequireAdminOrSupervisor(), pointReasonHandler.Update)
			pointReasons.DELETE("/:id", middleware.RequireAdminOrSupervisor(), pointReasonHandler.Delete)
		}

		points := api.Group("/points")
		{
			points.GET("", pointHandler.GetAll)
			points.GET("/student/:studentId", pointHandler.GetByStudentID)
			points.GET("/student/:studentId/summary", pointHandler.GetSummary)
			points.POST("", middleware.RequireAdminOrSupervisor(), pointHandler.GivePoint)
			points.POST("/bulk", middleware.RequireAdminOrSupervisor(), pointHandler.BulkGivePoints)
			points.PATCH("/:id/cancel", middleware.RequireAdminOrSupervisor(), pointHandler.Cancel)
			points.DELETE("/reset", middleware.RequireAdmin(), pointHandler.Reset)
		}

		duties := api.Group("/duties")
		{
			duties.GET("", dutyHandler.GetAll)
			duties.GET("/:id", dutyHandler.GetByID)
			duties.POST("", middleware.RequireAdminOrSupervisor(), dutyHandler.Create)
			duties.PUT("/:id", middleware.RequireAdminOrSupervisor(), dutyHandler.Update)
			duties.DELETE("/:id", middleware.RequireAdminOrSupervisor(), dutyHandler.Delete)
			duties.POST("/generate", middleware.RequireAdminOrSupervisor(), dutyHandler.Generate)
			duties.POST("/:id/swap", dutyHandler.Swap)
			duties.PATCH("/:id/complete", middleware.RequireAdminOrSupervisor(), dutyHandler.Complete)
		}

		api.GET("/audit-logs", middleware.RequireAdmin(), auditHandler.GetAll)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedAdmin(cfg *config.Config, authService *service.AuthService) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		return
	}
	_, err := authService.CreateUser(dto.CreateUserRequest{
		Email:    cfg.AdminEmail,
		Password: cfg.AdminPassword,
		Name:     "Admin",
		Role:     string(model.RoleAdmin),
	})
	if err != nil {
		log.Printf("Admin user may already exist: %v", err)
	}
}
