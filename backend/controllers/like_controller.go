package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context){
	articleId := ctx.Param("id")
	likekey := "article:" + articleId + ":likes"

	if err := global.Redisdb.Incr(likekey).Err(); err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{"message":"Successfully like"})
}

func GetLikes(ctx *gin.Context){
	articleId := ctx.Param("id")
	likekey := "article:" + articleId + ":likes"
	likes , err := global.Redisdb.Get(likekey).Result()
	// 修改逻辑，redis.Nil这边如果未命中才会返回，不是!=
	if err==redis.Nil{
		likes = "0"
	}else if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	ctx.JSON(http.StatusOK,gin.H{"Likes:":likes})
}