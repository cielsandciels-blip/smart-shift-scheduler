package database

import (
    "fmt"
    "log"
    "smart-shift-scheduler/internal/domain"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// NewDB: PostgreSQLに接続して、接続情報を返す
func NewDB() *gorm.DB {
    // (接続情報はそのまま)
    dsn := "host=127.0.0.1 user=manager password=manager123 dbname=shift_db port=5433 sslmode=disable TimeZone=Asia/Tokyo"

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("データベースへの接続に失敗しました:", err)
    }

    fmt.Println("データベース接続成功！")

    // ★ここを修正！ DailyRequirement を追加しました
    err = db.AutoMigrate(
        &domain.Staff{}, 
        &domain.Shift{}, 
        &domain.ShiftRequest{}, 
        &domain.DailyRequirement{}, // ★これを追加！
    )
    
    if err != nil {
        log.Fatal("テーブル作成に失敗しました:", err)
    }

    return db
}