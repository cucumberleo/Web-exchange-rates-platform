package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles"

func CreateArticle(ctx *gin.Context) {
	var article models.Article
	if err := ctx.ShouldBind(&article); err!=nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error":"Invalid import"})
		return
	}
	
	if err := global.Db.AutoMigrate(&article); err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err:=global.Redisdb.Del(cacheKey).Err();err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated,article)
}

func GetArticles(ctx *gin.Context){
	cachedata,err := global.Redisdb.Get(cacheKey).Result()
	var articles []models.Article
	// 如果redis缓存没有命中，需要从数据库中调入文章数据
	if err==redis.Nil{
		// 查询文章
		if err := global.Db.Find(&articles).Error; err != nil {
			if errors.Is(err,gorm.ErrRecordNotFound){
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error})
			}else{
				ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
			}
			return
		}
		// 把文章数据变成json数据存入redis
		articleJSON,err := json.Marshal(articles)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		if err:=global.Redisdb.Set(cacheKey,articleJSON,time.Minute*10).Err();err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,articles)
	}else if err!=nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}else{
		// 如果命中就把redis cache中的json数据解码成文章数据
		if err:= json.Unmarshal([]byte(cachedata),&articles);err!=nil{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		ctx.JSON(http.StatusOK,articles)
	}
}

// 通过id获取文章
func GetArticleByid(ctx *gin.Context){
	var article models.Article
	id := ctx.Param("id")

	if err := global.Db.Where("id=?",id).First(&article).Error; err!=nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}else{
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK,article)
}


