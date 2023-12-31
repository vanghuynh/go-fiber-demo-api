package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/golang-jwt/jwt/v5"
	"github.com/vanghuynh/fiber-api/database"
	"github.com/vanghuynh/fiber-api/routes"

	// docs are generated by Swag CLI, you have to import them.
	// replace with your own docs folder, usually "github.com/username/reponame/docs"
	_ "github.com/vanghuynh/fiber-api/docs"

	jwtware "github.com/gofiber/contrib/jwt"

	// need this lib to work with atlas command
	_ "ariga.io/atlas-go-sdk/recordriver"
	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/robfig/cron/v3"
)

var wg sync.WaitGroup

var c = make(chan int) // allocate a channel

// define a public and private key
// just for demo, in production you should use env
var PrivateKey *rsa.PrivateKey

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	app := fiber.New()

	// generate a public and private key pair
	var err error
	PrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		// stop application
		log.Fatalf("rsa.GenerateKey: %v", err)
	}

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	// user command to run migration, no need to integrate in code
	// connect to database
	database.ConnnectDb()

	// this route defined befoe JWT middleware is not authenticated
	app.Post("/api/user/login", LoginUser)

	// get list of trending coins
	// this api get data form coingecko, not from database
	app.Get("/api/coin", routes.GetTrendingCoins)

	// get list of all available coins
	// this api get data form coingecko, not from database
	app.Get("/api/coin/all", routes.GetCoinsList)

	// JWT middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    PrivateKey.Public(),
		},
	}))

	// all routes defined after JWT middleware must authenticate
	settupRoutes(app)

	// cron settup
	cron := cron.New()
	cron.AddFunc("@every 1m", func() {
		fmt.Println("start cron every minute")
		// coinList, err := routes.GetCoinGekoCoinsList()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println("coinList: ", coinList)
		// fmt.Println("coinList[0]: ", coinList[0])

	})
	cron.Start()
	fmt.Println("start cron")

	c1 := circle{
		radius: 10,
	}
	cp := c1
	fmt.Println("area: ", cp.area())

	// go routine
	fmt.Println("OS\t", runtime.GOOS)
	fmt.Println("ARCH\t", runtime.GOARCH)
	fmt.Println("CPUs\t", runtime.NumCPU())
	fmt.Println("Goroutines\t", runtime.NumGoroutine())
	// make function as go routine
	go taskOne()
	<-c // wait for taskOne finish
	taskTwo()

	app.Listen(":3000")

}

type circle struct {
	radius float64
}

func (c *circle) area() float64 {
	return 3.14 * c.radius * c.radius
}

func taskOne() {
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)
		fmt.Println("taskOne")
	}
	c <- 1 // send a signal, value does not matter
}

func taskTwo() {
	for i := 0; i < 10; i++ {
		time.Sleep(6 * time.Second)
		fmt.Println("taskTwo")
	}
}

func settupRoutes(app *fiber.App) {
	app.Post("/api/user", routes.CreateUser)
	app.Get("/api/user", routes.GetUsers)
	app.Get("/api/user/:id", routes.GetUser)
	app.Put("/api/user/:id", routes.UpdateUser)
	app.Delete("/api/user/:id", routes.DeleteUser)
	// app.Post("/api/user/login", LoginUser)

	app.Post("/api/product", routes.CreateProduct)
	app.Get("/api/product", routes.GetProducts)

	app.Post("/api/order", routes.CreateOrder)
	app.Get("/api/order", routes.GetOrders)
}

// LoginUser godoc
// @Summary      Login user
// @Description  login user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        data body UserLoginRequestDto true "Login input"
// @Success      200  {object}  UserLoginResponseDto
// @Router       /api/user/login [post]
func LoginUser(c *fiber.Ctx) error {
	// get input from request
	var userLoginRequestDto UserLoginRequestDto
	if err := c.BodyParser(&userLoginRequestDto); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	// generate token
	if userLoginRequestDto.UserName != "admin" || userLoginRequestDto.Password != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON("Invalid username or password")
	}
	// create claims
	claims := jwt.MapClaims{
		"username": userLoginRequestDto.UserName,
		"password": userLoginRequestDto.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// generate token
	tokenString, err := token.SignedString(PrivateKey)
	if err != nil {
		log.Printf("token.SignedString: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var userLoginResponseDto = UserLoginResponseDto{
		Token: tokenString,
	}
	return c.Status(200).JSON(userLoginResponseDto)
}

type UserLoginRequestDto struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type UserLoginResponseDto struct {
	Token string `json:"token"`
}
