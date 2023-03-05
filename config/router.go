package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"promptscroll/controller"
	"promptscroll/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 路由中间件
	r.Use(middleware.ErrorMiddleware())

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth 路由组
	authRouter := r.Group("/api/auth")
	{
		authController := controller.NewAuthController()
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/register", authController.Register)
	}

	// User 路由组
	userRouter := r.Group("/api/user")
	userRouter.Use(middleware.AuthMiddleware())
	{
		userController := controller.NewUserController()
		userRouter.GET("/:id", userController.GetUserByID)
	}

	// 404 Not Found 路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found"})
	})

	return r
}
