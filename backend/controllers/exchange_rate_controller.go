package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate
	// 将请求头里面的json数据填进结构体实例
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exchangeRate.Date = time.Now()
	// 建表过程
	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, exchangeRate)
}

func GetExchangeRates(ctx *gin.Context) {
	var exchangeRates []models.ExchangeRate

	if err := global.Db.Find(&exchangeRates).Error; err != nil {
		// 如果确实没找到记录
		if errors.Is(err,gorm.ErrRecordNotFound){
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error})
		}else{
			// 数据有但是找不到
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":err})
		}
		return
	}
	ctx.JSON(http.StatusOK, exchangeRates)
}
