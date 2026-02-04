package database

import (
	"fmt"
	"log"
	"smart-shift-scheduler/internal/domain" // domainパッケージのstructを使うため

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB: PostgreSQLに接続して、接続情報を返す
func NewDB() *gorm.DB {
	// Dockerで設定したユーザー名やパスワード
	// 変更前
// dsn := "host=127.0.0.1 user=user password=password dbname=shift_db port=5432 ..."

// 変更後 (ポート5433, user=manager, pass=manager123)
dsn := "host=127.0.0.1 user=manager password=manager123 dbname=shift_db port=5433 sslmode=disable TimeZone=Asia/Tokyo"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}

	fmt.Println("データベース接続成功！")

	// ★マイグレーション（超重要）
	// Goのstruct（Staffなど）を見て、自動でSQLのテーブルを作ってくれる機能
	err = db.AutoMigrate(&domain.Staff{}, &domain.Shift{}, &domain.ShiftRequest{})
	if err != nil {
		log.Fatal("テーブル作成に失敗しました:", err)
	}

	return db
}