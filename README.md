## 后端功能说明

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
Use `global.Db.AutoMigrate(&user)` to update the database with new information
```go
if err := global.Db.AutoMigrate(&user); err != nil {
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
