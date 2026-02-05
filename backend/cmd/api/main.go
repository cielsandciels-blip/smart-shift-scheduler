package main

import (
	"fmt"
	"path/filepath"
	"smart-shift-scheduler/internal/handler"
	"smart-shift-scheduler/internal/infrastructure/database"
	"smart-shift-scheduler/internal/infrastructure/engine"
	"smart-shift-scheduler/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.NewDB()
	scriptPath := filepath.Join("..", "engine", "main.py")
	shiftEngine := engine.NewShiftEngine(scriptPath)

	// Staff
	staffRepo := database.NewStaffRepository(db)
	staffUsecase := usecase.NewStaffUsecase(staffRepo)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	// Shift & Request & Requirement (★ここを拡張)
	shiftRepo := database.NewShiftRepository(db)
	requestRepo := database.NewRequestRepository(db)
	requireRepo := database.NewRequirementRepository(db) // ★追加1: 必要人数の保存場所
	
	// ★追加2: 引数が5つになりました (engine, staffRepo, shiftRepo, requestRepo, requireRepo)
	shiftUsecase := usecase.NewShiftUsecase(shiftEngine, staffRepo, shiftRepo, requestRepo, requireRepo)
	
	shiftHandler := handler.NewShiftHandler(shiftUsecase)
	requestHandler := handler.NewRequestHandler(shiftUsecase)

	r := gin.Default()
	r.Static("/web", "../frontend")

	api := r.Group("/api")
	{
		api.POST("/staff", staffHandler.Create)
		api.GET("/staff", staffHandler.List)
		api.DELETE("/staff/:id", staffHandler.Delete)
		
		api.POST("/shift", shiftHandler.Generate)
		api.GET("/shift", shiftHandler.List)
		api.PUT("/shift/:id", shiftHandler.Update)
		api.DELETE("/shift/:id", shiftHandler.Delete)

		api.POST("/request", requestHandler.Create)
		api.GET("/request", requestHandler.List)
		api.DELETE("/request/:id", requestHandler.Delete)

		// ★追加3: 必要人数設定のAPI
		api.POST("/requirement", shiftHandler.SaveRequirement)
		api.GET("/requirement", shiftHandler.ListRequirements)
		api.DELETE("/requirement/:id", shiftHandler.DeleteRequirement) // 追加
	
		api.GET("/export", shiftHandler.Export)
	}

	fmt.Println("サーバーを起動します... http://localhost:8080/web/index.html")
	r.Run(":8080")
}