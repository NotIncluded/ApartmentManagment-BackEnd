package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/PunMung-66/ApartmentSys/config"
	"github.com/PunMung-66/ApartmentSys/controller"
	"github.com/PunMung-66/ApartmentSys/internal/auth"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := flag.String("e", "dev", "environment (local|dev)")
	setup := flag.Bool("setup", false, "setup database")
	flag.Parse()

	envFile := ".env." + *env
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Println("Error loading env file")
	}

	port := os.Getenv("PORT")
	secret := os.Getenv("JWT_SECRET")

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		res := map[string]any{
			"status":  http.StatusOK,
			"version": "0.0.1",
		}
		c.JSON(http.StatusOK, res)
	})

	db, err := config.ConnectDatabase()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	if *setup {
		fmt.Println("Running initial setup (AutoMigrate)...")

		db.AutoMigrate(
			&model.User{},
			&model.Room{},
			&model.Contract{},
			&model.UtilityRate{},
			&model.UtilityUsage{},
			&model.Bill{},
			&model.Payment{},
		)

		fmt.Println("Setup completed!")
		return
	}
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userRoute := r.Group("/user")
	userRoute.Use(auth.Protect([]byte(secret), "ADMIN"))
	{
		userRoute.POST("/create", userController.CreateUser)
		userRoute.DELETE("/:id", userController.DeleteUser)
	}

	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService, []byte(secret))

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", authController.LoginHandler)
		authRoute.POST("/register", authController.RegisterHandler)
	}

	r.Run(":" + port)
}
