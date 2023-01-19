package models

import "gorm.io/gorm"

type User struct{
    gorm.Model
    Username string `gorm:"unique"json:"username"`
    Email string `json:"email"`
    Age string `json:"age"`
    Gender string `json:"gender"`
    Height string `json:"height"`
    Weight string `json:"weight"`
    Goal string `json:"goal"`
    Activity string `json:"activity"`
    Password string `json:"password"`
    Food []Food `json:"food"`
}

type Food struct{
    gorm.Model
    Name string `gorm:""json:"name"`
    Calories string `json:"calories"`
    Ingredient []Ingredient `json:"ingredients"`
    UserId uint `json:"userId"`
}

type Ingredient struct{
    gorm.Model
    Name string `gorm:""json:"name"`
    FoodId uint `json:"foodId"`
}