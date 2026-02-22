package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/grahagandangr/nexboard-be/config"
	"github.com/grahagandangr/nexboard-be/handlers"
	"github.com/grahagandangr/nexboard-be/middleware"
	"github.com/grahagandangr/nexboard-be/repositories"
	"github.com/grahagandangr/nexboard-be/services"
)

func main() {
	// 1. Load configuration
	config.LoadConfig()

	// 2. Connect to database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	// 3. Initialize repositories
	userRepo := repositories.NewUserRepository(config.DB)
	workspaceRepo := repositories.NewWorkspaceRepository(config.DB)
	boardRepo := repositories.NewBoardRepository(config.DB)
	statusRepo := repositories.NewStatusRepository(config.DB)
	taskRepo := repositories.NewTaskRepository(config.DB)

	// 4. Initialize services
	authService := services.NewAuthService(userRepo)
	workspaceService := services.NewWorkspaceService(workspaceRepo, userRepo)
	boardService := services.NewBoardService(boardRepo, workspaceRepo, userRepo)
	statusService := services.NewStatusService(statusRepo)
	taskService := services.NewTaskService(taskRepo, boardRepo, statusRepo, userRepo, workspaceRepo)

	// 5. Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	workspaceHandler := handlers.NewWorkspaceHandler(workspaceService)
	boardHandler := handlers.NewBoardHandler(boardService)
	statusHandler := handlers.NewStatusHandler(statusService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// 6. Setup Gin router
	router := gin.Default()

	// 7. Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "NexBoard API is running smoothly."})
	})

	// 8. API routes setup
	api := router.Group("/api")
	{
		// ---- PUBLIC ROUTES ----
		users := api.Group("/users")
		{
			users.POST("/login", authHandler.Login)
			users.POST("/register", authHandler.Register)
		}

		// ---- PROTECTED ROUTES ----
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			// User Profile
			protected.GET("/users/profile", authHandler.GetProfile)
			protected.PUT("/users/profile", authHandler.UpdateProfile)

			// Workspaces
			workspaces := protected.Group("/workspaces")
			{
				workspaces.POST("", workspaceHandler.CreateWorkspace)
				workspaces.GET("", workspaceHandler.GetUserWorkspaces)
				workspaces.GET("/:external_id", workspaceHandler.GetWorkspace)
				workspaces.PUT("/:external_id", workspaceHandler.UpdateWorkspace)
				workspaces.DELETE("/:external_id", workspaceHandler.DeleteWorkspace)

				// Workspace Members
				members := workspaces.Group("/:external_id/members")
				{
					members.GET("", workspaceHandler.GetMembers)
					members.POST("", workspaceHandler.InviteMember)
					members.PUT("/:user_ext_id", workspaceHandler.UpdateMemberRole)
					members.DELETE("/:user_ext_id", workspaceHandler.RemoveMember)
				}

				// Workspace Boards
				boards := workspaces.Group("/:external_id/boards")
				{
					boards.POST("", boardHandler.CreateWorkspaceBoard)
					boards.GET("", boardHandler.GetWorkspaceBoards)
				}
			}

			// Boards (direct manipulation)
			boards := protected.Group("/boards")
			{
				boards.GET("/:external_id", boardHandler.GetBoard)
				boards.PUT("/:external_id", boardHandler.UpdateBoard)
				boards.DELETE("/:external_id", boardHandler.DeleteBoard)

				// Board Tasks
				tasks := boards.Group("/:external_id/tasks")
				{
					tasks.POST("", taskHandler.CreateBoardTask)
					tasks.GET("", taskHandler.GetBoardTasks)
				}
			}

			// Tasks (direct manipulation)
			tasks := protected.Group("/tasks")
			{
				tasks.PUT("/:external_id", taskHandler.UpdateTask)
				tasks.DELETE("/:external_id", taskHandler.DeleteTask)
				tasks.PATCH("/:external_id/status", taskHandler.MoveTask)
				tasks.PATCH("/:external_id/assign", taskHandler.AssignTask)
			}

			// Status Master Data
			statuses := protected.Group("/statuses")
			{
				statuses.POST("", statusHandler.CreateStatus)
				statuses.GET("", statusHandler.GetAllStatuses)
				statuses.GET("/:external_id", statusHandler.GetStatus)
				statuses.PUT("/:external_id", statusHandler.UpdateStatus)
				statuses.DELETE("/:external_id", statusHandler.DeleteStatus)
			}
		}
	}

	// 9. Setup graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		config.CloseDatabase()
		os.Exit(0)
	}()

	// 10. Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
