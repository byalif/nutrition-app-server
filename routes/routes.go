package routes

import (
        "github.com/gofiber/fiber/v2"
        "github.com/byalif/server/controllers"
)

func Setup(app * fiber.App) {
    app.Post("/createUser", controllers.Register)
    app.Post("/getUser", controllers.GetUser)
    app.Post("/login", controllers.Login)
    app.Post("/addFood", controllers.AddFood)
    app.Post("/deleteCookie", controllers.DeleteCookie)
    app.Get("/removeFood/:foodId", controllers.RemoveFood)
    app.Get("/filter/:search/:userId", controllers.SearchFilters)
}


