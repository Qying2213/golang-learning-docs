package main

import (
	"fmt"
	"log"
	"test1/config"
	"test1/handler"
	"test1/middleware"
	"test1/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ========== 1. 初始化数据库 ==========
	db, err := initDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功")

	// ========== 2. 创建 Handler ==========
	userHandler := handler.NewUserHandler(db)

	// ========== 3. 创建 Gin 引擎 ==========
	r := gin.Default()

	// ========== 4. 注册中间件 ==========
	// CORS 中间件（全局）
	r.Use(middleware.CORS())

	// ========== 5. 注册路由 ==========
	// TODO: 任务5 - 注册路由

	// 公开接口（不需要登录）
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
	}

	// 需要登录的接口
	// TODO: 添加 JWT 中间件
	auth := r.Group("/api")
	auth.Use(middleware.JWT()) // JWT 中间件
	{
		// TODO: 注册需要登录的路由
		// auth.GET("/user/profile", userHandler.GetProfile)
		// auth.PUT("/user/profile", userHandler.UpdateProfile)
		// auth.PUT("/user/password", userHandler.UpdatePassword)
		// auth.GET("/users", userHandler.ListUsers)
	}

	// ========== 6. 启动服务器 ==========
	fmt.Println("服务器启动在", config.ServerPort)
	r.Run(config.ServerPort)
}

// initDB 初始化数据库连接
func initDB() (*gorm.DB, error) {
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移（创建表）
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
