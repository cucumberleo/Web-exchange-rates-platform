package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LikeArticle(ctx *gin.Context){
	articleId := ctx.Param("id")
	name,exists := ctx.Get("username")
	if !exists{
		ctx.JSON(http.StatusNotFound,gin.H{"error":"Do not found the user"})
		return
	}
	// 修改键逻辑
	likekey := "article:" + articleId + ":likers"
	// 修改逻辑，每个用户只能点一次赞，用redis set
	isMember , err := global.Redisdb.SIsMember(likekey,name).Result()
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check like status"})
		return
	}
	if isMember{
		if err := global.Redisdb.SRem(likekey,name).Err(); err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Could not cancel like"})
			return
		}
	}else{
		if err := global.Redisdb.SAdd(likekey,name).Err(); err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":"Could not like"})
			return
		}
	}
	newLikesCount, err := global.Redisdb.SCard(likekey).Result()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取更新后的点赞数 (Could not retrieve updated like count)"})
        return
    }
	// 返回操作成功信息、最新的点赞总数以及用户最新的点赞状态
	ctx.JSON(http.StatusOK, gin.H{
        "message": "操作成功 (Operation successful)",
        "likes": newLikesCount,
        "is_liked": !isMember, // 最新的状态与之前的状态相反
    })
}
func GetLikes(ctx *gin.Context){
	articleId := ctx.Param("id")
	likekey := "article:" + articleId + ":likers"
	likes, err := global.Redisdb.SCard(likekey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}


func Islike(ctx *gin.Context){
	articleId := ctx.Param("id")
	name, exists := ctx.Get("username")
	if !exists {
		// 如果用户未登录，那么他肯定没有点赞
		ctx.JSON(http.StatusOK, gin.H{"is_liked": false})
		return
	}
	likekey := "article:" + articleId + ":likers"
	isMember , err := global.Redisdb.SIsMember(likekey,name).Result()
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"is_liked":isMember})
}