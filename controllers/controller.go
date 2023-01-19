package controllers

import (
     "github.com/gofiber/fiber/v2"
     "github.com/byalif/server/models"
    "github.com/byalif/server/config"
    "time"
    "log"
    "strings"
     "golang.org/x/crypto/bcrypt"
     "strconv"
     "github.com/dgrijalva/jwt-go"
)

const SecretKey = "secret"

type MyCustomClaims struct {
	UserId int `json:"userid"`
	jwt.StandardClaims
}

func AddFood(c *fiber.Ctx) error {
    var data map[string]string


    if err := c.BodyParser(&data); err!= nil {
        return err
    }

    arr := strings.Fields(data["ingredients"])

    userId, _ := strconv.Atoi(data["userId"])
    calories, _ := strconv.Atoi(data["calories"])

    food := models.Food{
        Name: data["name"],
        Calories: calories,
        UserId: uint(userId),
    }

    config.DB.Create(&food)

    for _, v := range arr {
        go createIngredient(v, food.ID)
    }
    return c.JSON(food)
}

func RemoveFood(c *fiber.Ctx) error {
    id := c.Params("foodId")

    config.DB.Delete(&models.Food{}, id)

    return c.JSON(fiber.Map{
                "status" : "deleted",
    })

}

func DeleteCookie(c *fiber.Ctx) error {
    cookie := fiber.Cookie{
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-3 * time.Hour),
        HTTPOnly: false,
    }

    c.Cookie(&cookie)

    return c.JSON(fiber.Map{
                "status" : "deleted",})

}

func SearchFilters(c *fiber.Ctx) error {
    search := c.Params("search")
    id := c.Params("userId")
    

    search = strings.ToLower(search)

    var foods []models.Food
    var response []models.Food


    config.DB.Where("user_id=?", id).Preload("Ingredient").Find(&foods)

    if search == "6910" {
      return c.JSON(foods)
    }


    for _, v := range foods {

        var str = v.Name
        str = strings.ToLower(str)
        if strings.Contains(str, search) {
            response = append(response, v)
            continue;
        }

        var ingredients []models.Ingredient
        config.DB.Where("food_id=?", v.ID).Find(&ingredients)
        for _, h := range ingredients {
            var str = h.Name
            str = strings.ToLower(str)
            if strings.Contains(str, search) {
                response = append(response, v)
                break;
            }
        }
    }


    return c.JSON(response)
}


func createIngredient(v string, Id uint){
        ingredient := models.Ingredient{
                Name: v,
                FoodId: Id,
        }
        config.DB.Create(&ingredient)
}


func Register(c *fiber.Ctx) error {
    var data map[string]string

    if err := c.BodyParser(&data); err!= nil {
        return err
    }

    height, _ := strconv.Atoi(data["height"])
    weight, _ := strconv.Atoi(data["weight"])
    age, _ := strconv.Atoi(data["age"])

    password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14);

    user := models.User{
        Username: data["username"],
        Email: data["email"],
        Password: string(password),
        Age: age,
        Gender: data["gender"],
        Activity: data["activity"],
        Goal: data["goal"],
        Height: height,
        Weight: weight,
    }

    config.DB.Create(&user)

    return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error{
    cookie := c.Cookies("jwt")


    token, err := jwt.ParseWithClaims(string(cookie), &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(SecretKey), nil
    })

    if err!= nil {
          c.Status(fiber.StatusUnauthorized)
            return c.JSON(fiber.Map{
                "error" : "UNAUTHORIZED",
            })
    }
    claims := token.Claims.(*MyCustomClaims)

    user := models.User{}
    config.DB.Where("ID=?", claims.UserId).Preload("Food").Preload("Food.Ingredient").Find(&user)
    return c.JSON(user)
}

func Login(c *fiber.Ctx) error{
    var body map[string]string

    if err := c.BodyParser(&body); err!= nil {
        log.Fatal(err)
    }

    user:= models.User{ ID: -1}

    config.DB.Where("username=?", body["username"]).Preload("Food").Preload("Food.Ingredient").First(&user)
    if user.ID == -1 {
            // c.Status(fiber.StatusNotFound)
            return c.JSON(user)
    }else {
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password),  []byte(body["password"])); err != nil {
            c.Status(fiber.StatusUnauthorized)
            return c.JSON(fiber.Map{
                "error" : "WRONG_PASSWORD",
            })
        }
    }

    claims := MyCustomClaims{
        int(user.ID),
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
            Issuer:    "test",
        },
    }
    

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(SecretKey))

    c.Request().Header.Set("token",tokenString)

    if err!= nil {
         c.Status(fiber.StatusInternalServerError)
            return c.JSON(fiber.Map{
                "error" : "COULD_NOT_LOG_IN",
            })
    }

    cookie:= fiber.Cookie{
        Name: "jwt",
        Value: tokenString,
        Expires: time.Now().Add(time.Hour * 24),
        HTTPOnly: false,
    }

    c.Cookie(&cookie)
    c.Status(fiber.StatusOK)
    return c.JSON(user)
}