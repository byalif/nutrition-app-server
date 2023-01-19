package main

import (
    "os"
    "log"
    "github.com/byalif/server/config"
    "github.com/byalif/server/routes"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func main() {
    config.Connect()

    app := fiber.New()

    app.Use(cors.New(cors.Config{
           AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
        AllowOrigins:     "https://cal-tracker.herokuapp.com",
        AllowCredentials: true,
        AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
    }))

    routes.Setup(app)
    
    dotenv := goDotEnvVariable("PORT")

    str := ":"+ dotenv;
    app.Listen(str)
}