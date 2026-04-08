package database

import (
	"fmt"
	"log"

	"silverlord/backend/config"
	"silverlord/backend/models"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var DB *xorm.Engine

func Init(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	engine, err := xorm.NewEngine("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}

	if err := engine.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	// 自動建表
	if err := engine.Sync2(
		new(models.Product),
		new(models.SpotPrice),
		new(models.Order),
		new(models.OrderItem),
		new(models.SystemConfig),
	); err != nil {
		log.Fatalf("DB sync failed: %v", err)
	}

	// 確保 SystemConfig 有預設資料（id=1）
	count, _ := engine.Count(new(models.SystemConfig))
	if count == 0 {
		engine.Insert(&models.SystemConfig{Id: 1})
	}

	DB = engine
	log.Println("Database connected and synced.")
}
