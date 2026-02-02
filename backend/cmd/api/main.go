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
	// 1. DB接続
	db := database.NewDB()

	// 2. Pythonエンジンの準備
	scriptPath := filepath.Join("..", "engine", "main.py")
	shiftEngine := engine.NewShiftEngine(scriptPath)

	// 3. 依存関係の組み立て
	
	// スタッフ管理
	staffRepo := database.NewStaffRepository(db)
	staffUsecase := usecase.NewStaffUsecase(staffRepo)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	// シフト作成
	shiftRepo := database.NewShiftRepository(db) // ★追加
	// ★修正点: 引数を3つ渡す (engine, staffRepo, shiftRepo)
	shiftUsecase := usecase.NewShiftUsecase(shiftEngine, staffRepo, shiftRepo) 
	shiftHandler := handler.NewShiftHandler(shiftUsecase)

	// 4. Webサーバーの設定
	r := gin.Default()
	r.Static("/web", "../frontend") // フロントエンド配信

	api := r.Group("/api")
	{
		api.POST("/staff", staffHandler.Create)
		api.GET("/staff", staffHandler.List)
		
		api.POST("/shift", shiftHandler.Generate)
		api.GET("/shift", shiftHandler.List)

		// 追加: 更新と削除
		api.PUT("/shift/:id", shiftHandler.Update)
		api.DELETE("/shift/:id", shiftHandler.Delete)
	}

	fmt.Println("サーバーを起動します... http://localhost:8080/web/index.html")
	r.Run(":8080")
}