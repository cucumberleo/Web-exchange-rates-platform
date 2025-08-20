## Backend Description

### Register function
   **1. Bind JSON information with user struct**
   `user.go`
```go
    type User struct {
        gorm.Model
        // 保证用户名不重复
        Username string `gorm:"unique"`
         // Password
        Password string
    }
```
   **gorm.Model.go**:proivde a convenient struct to operate db functions
```go
    type Model struct {
        ID        uint `gorm:"primarykey"`
        CreatedAt time.Time
        UpdatedAt time.Time
        DeletedAt DeletedAt `gorm:"index"`
    }
```
`auth_controller.go(part)`
```go
var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
```
   **2. Encription methods**
   **Part 1 : Hash password**
   use `bcrypt.GenerateFromPassword([]byte(passward),cost)` to create hash password
```go
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}
```
**Part 2 : Digital signiture**
Commonly, the process of digital signiture can be divided into 3 parts:
1. `Header`: The type of token + signing algorithms claims;  use `jwt.NewWithClaims()`
2. `Payload`: Including Claims and statement of user and other data; use `jwt.MapClaims` to claim `username` and `expiration`
3. `Sign token`: Use `secret key` to generate `signedToken`

```go
func GenerateJWT(username string) (string, error) {
	// Header: 令牌类型JMT + 签名算法HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Payload：包含声明Claims，关于用户和其他数据的陈述
		"username": username,
		// 有效期设为3天
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	// secret是密钥
	signedToken, err := token.SignedString([]byte("secret"))
	// Bearer是一种认证方式
	return "Bearer " + signedToken, err
}
```
**3. Write data into Database**
Use `global.Db.Create(&user)` to update the database with new information
```go
if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
```

**4. Config the router setting**
1. Use `router.Group()` to set router group
2. Use `auth.POST()` to post register function package
```go
auth := r.Group("/api/auth")
{
	auth.POST("/register", controllers.Register)
}
```

### Login function
**1. Define input format struct and Bindjson**
```go
var input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
...
```

**2. Use gorm.db function `global.Db.Where("username = ?",input.Username).First(&user)` to find the first record of username that matched**
```go
if err := global.Db.Where("username = ?", input.Username).First((&user)).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials1"})
		return
	}
```
**3. Use `bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))` to check if the corresponding password is correct**
if not matched , return 401 error
```go
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```


**4. Generate token and return**
Also use `GenerateJWT(user.Username)` as used in Register func.


**5. router setting**
```go
auth := r.Group("/api/auth")
{
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register)
}
```

### Upload Exchangerate func
**1. What message should we include in exchange rates? Is there any relationship between each piece of information?**
```go
type ExchangeRate struct {
	ID           uint      `gorm:"primarykey" json:"_id"`
	FromCurrency string    `json:"fromCurrency" binding:"required"`
	ToCurrency   string    `json:"toCurrency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
}
```
`FromCurrency` : Your currency
`ToCurrency` : Exchanged currency
`Rate` : Exchange rate
These three information should be binded together to make calculation.


**2.Two key functions we should write in `exchange_rate_controller.go`**
##### 1.`func CreateExchangeRate(ctx *gin.Context)`:Create a record of exchangerate
·Thinking: The information we upload is `Fromcurrency`,`Tocurrency`,`rate`; Besides,Id is set as primarykey and increase automatically. But our message struct include also date, we should not forget to manually set the date information.
**Bind json information**
```go
if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
}
exchangeRate.Date = time.Now()
```

**Write information into DB**
```go
// Try to fit the format
if err:=global.Db.AutoMigrate(&exchangeRate); err!=nil{
	// 500 error
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
}
// Create the table record in sql
if err:=global.Db.Create(&exchangeRate); err!=nil{
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
}
```

##### 2. func GetExchangeRates(ctx *gin.Context): Get the exchange rate
Thinking: As we do not want to return all the message in exchangerate struct, so we only need to instantiate slice of the struct.`[]models.ExchangeRate`
Use `errors.Is` to judge if the error is the certain error
```go
if err:= global.Db.Find(&exchangeRates).Error; err!=nil{
	// Case1:The record is acturally not in the Db
	// Case2:The record is in Db but you can not fetch it
	if errors.Is(err,gorm.ErrorRecordNotFound){
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error})
	}else{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Err})
	}
	return
	
}
```
**3.Setup router**
Thinking: The exchangeRate message is open to public in my design, but I want some later functions and create exchange rate are only available to user that login,how should I achieve this?  Answer: Use middleware!

##### AuthMiddleware
Use `routergroup.Use(middleware HandlerFunc)` to verify login status.
Offical code：
```go
func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}
```




