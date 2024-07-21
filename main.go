package main

import (
    "math/rand"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "log"
    "fmt"
    "time"
)
var db *gorm.DB
const shortCodeLength = 8
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


type User struct {
    ID        uint      `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Email     string    `gorm:"type:varchar(100);uniqueIndex"`
    Password  string
}

type URL struct {
    gorm.Model
    ShortCode string `gorm:"uniqueIndex"`
    OriginalURL string
    UserID uint
}

func initDB() {
    dsn := "root:abulata@2004@tcp(127.0.0.1:3306)/shortener?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v\n", err)
    } else {
        fmt.Println("Connected to database successfully")
    }

    
    db.AutoMigrate(&User{})
}

func register(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = string(hashedPassword)

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func login(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var dbUser User
    if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func generateShortCode() string {
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)
    shortCode := make([]byte, shortCodeLength)
    for i := range shortCode {
        shortCode[i] = charset[random.Intn(len(charset))]
    }
    return string(shortCode)
}

func shortenURL(c *gin.Context) {
    var request struct {
        OriginalURL string `json:"original_url"`
    }
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    shortCode := generateShortCode()
    url := URL{
        ShortCode: shortCode,
        OriginalURL: request.OriginalURL,
        // Assume user ID is 1 for simplicity; replace with actual user ID
        UserID: 1,
    }

    if err := db.Create(&url).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"short_code": shortCode})
}

func getOriginalURL(c *gin.Context) {
    shortCode := c.Param("short_code")
    var url URL
    if err := db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }

    c.Redirect(http.StatusFound, url.OriginalURL)
}

func main() {
    initDB()
    r := gin.Default()
    r.POST("/register", register)
    r.POST("/login", login)

    r.POST("/shorten", shortenURL)
    r.GET("/:short_code", getOriginalURL)
    
    r.Run(":8080")
}
