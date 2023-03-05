package main

import (
	"log"

	"promptscroll/config"
	"promptscroll/router"
)

func main() {
	// 初始化配置和路由
	cfg := config.LoadConfig()
	r := router.SetupRouter(cfg)

	// 開啟資料庫連線
	db := config.OpenDatabase(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close database: %s", err)
		}
	}()

	// 啟動 HTTP 伺服器
	if err := r.Run(cfg.Port); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
