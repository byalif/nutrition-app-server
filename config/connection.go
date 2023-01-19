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
   
    connection, err := gorm.Open(mysql.Open("bb78eb07479e92:b8dd168b@tcp(us-cdbr-east-06.cleardb.net:3306)/heroku_0d3f168b737c65d?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
    if err!= nil {
        panic(err)
    }

    DB = connection

    connection.AutoMigrate(&models.User{})
    connection.AutoMigrate(&models.Food{})
    connection.AutoMigrate(&models.Ingredient{})
}