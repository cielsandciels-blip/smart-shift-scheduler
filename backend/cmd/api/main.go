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

	// Shift & Request (★ここが変わりました)
	shiftRepo := database.NewShiftRepository(db)
	requestRepo := database.NewRequestRepository(db) // 追加
	
	// 引数が4つになりました (engine, staffRepo, shiftRepo, requestRepo)
	shiftUsecase := usecase.NewShiftUsecase(shiftEngine, staffRepo, shiftRepo, requestRepo)
	
	shiftHandler := handler.NewShiftHandler(shiftUsecase)
	requestHandler := handler.NewRequestHandler(shiftUsecase) // 追加

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

		// ★リクエスト用のAPIを追加
		api.POST("/request", requestHandler.Create)
		api.GET("/request", requestHandler.List)
		api.DELETE("/request/:id", requestHandler.Delete)
	}

	fmt.Println("サーバーを起動します... http://localhost:8080/web/index.html")
	r.Run(":8080")
}