package config

import (
    "gorm.io/gorm"
    "github.com/byalif/server/models"
    "log"
    "gorm.io/driver/mysql"
)

var  DB *gorm.DB

func Connect(){
    log.Println("Connection established")
   
    connection, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/goProj?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
    if err!= nil {
        panic(err)
    }

    DB = connection

    connection.AutoMigrate(&models.User{})
    connection.AutoMigrate(&models.Food{})
    connection.AutoMigrate(&models.Ingredient{})
}