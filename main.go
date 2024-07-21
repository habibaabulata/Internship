package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
    "log"
)

var db *gorm.DB

func initDB() {
    dsn := "root:abulata@2004@tcp(127.0.0.1:3306)/urlshortener?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }
}

func main() {
    initDB()
    r := gin.Default()
    // Add your routes here
    r.Run(":8080")
}
