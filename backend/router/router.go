package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST","OPTIONS"},
    AllowHeaders:     []string{"Origin","Content-Type","Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))
	// 设置路由组
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}
	api := r.Group("/api")
	// 不需要注册可以获取汇率
	api.GET("/exchangeRates", controllers.GetExchangeRates)
	// 括号内的都是需要注册才能做的事情
	api.Use(middlewares.Authmiddleware())
	{
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		api.POST("/articles",controllers.CreateArticle)
		api.GET("articles",controllers.GetArticles)
		api.GET("/articles/:id",controllers.GetArticleByid)

		api.POST("/articles/:id/like",controllers.LikeArticle)
		api.GET("articles/:id/like",controllers.GetLikes)
	}
	return r
}
