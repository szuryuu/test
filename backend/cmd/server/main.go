package main

import (
	"fmt"
	"log/slog"
	"os"

	"kasiraiai/backend/config"
	"kasiraiai/backend/internal/handler"
	"kasiraiai/backend/internal/middleware"
	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/internal/service"
	"kasiraiai/backend/pkg/ai"
	"kasiraiai/backend/pkg/fonnte"

	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("gagal memuat konfigurasi", "error", err)
		os.Exit(1)
	}

	// Database
	pool, err := repository.InitDB(cfg)
	if err != nil {
		slog.Error("gagal koneksi database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Auto-run migrations
	if err := repository.RunMigrations(cfg); err != nil {
		slog.Error("gagal menjalankan migration", "error", err)
		os.Exit(1)
	}

	// Repositories
	umkmRepo := repository.NewUmkmRepo(pool)
	transactionRepo := repository.NewTransactionRepo(pool)

	// External clients
	aiClient := ai.NewClient(cfg)
	aiParser := ai.NewParser(aiClient)
	fonnteClient := fonnte.NewClient(cfg)

	// Services
	authService := service.NewAuthService(umkmRepo, cfg)
	dashboardService := service.NewDashboardService(transactionRepo)
	kurService := service.NewKurService(transactionRepo)
	reportService := service.NewReportService(transactionRepo, fonnteClient)
	transactionService := service.NewTransactionService(transactionRepo, aiParser, fonnteClient, dashboardService, reportService, kurService)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	webhookHandler := handler.NewWebhookHandler(transactionService, umkmRepo, fonnteClient)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	umkmHandler := handler.NewUmkmHandler(umkmRepo)

	// Router
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))

	// Health check — no auth, before api group
	r.GET("/health", handler.HealthCheck)
	r.HEAD("/health", handler.HealthCheck)

	// Webhook — no auth
	r.POST("/webhook/whatsapp", webhookHandler.Handle)
	r.GET("/webhook/whatsapp", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := r.Group("/api/v1")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		// KUR & Report handlers (thin wrappers over services)
		kurHandler := handler.NewKurHandler(kurService)
		reportHandler := handler.NewReportHandler(reportService, umkmRepo)

		protected := api.Group("")
		protected.Use(middleware.Auth(authService))
		{
			protected.GET("/transactions", transactionHandler.List)
			protected.POST("/transactions", transactionHandler.Create)
			protected.DELETE("/transactions/:id", transactionHandler.Delete)

			protected.GET("/dashboard/summary", dashboardHandler.Summary)
			protected.GET("/dashboard/categories", dashboardHandler.Categories)

			protected.GET("/kur/score", kurHandler.GetScore)
			protected.POST("/kur/recalculate", kurHandler.Recalculate)

			protected.POST("/reports/monthly", reportHandler.GenerateMonthly)

			protected.GET("/umkm/profile", umkmHandler.GetProfile)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	handler.IsReady.Store(true)
	slog.Info("server berjalan", "addr", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("server gagal", "error", err)
		os.Exit(1)
	}
}
