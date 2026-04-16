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
	"github.com/gin-contrib/cors"
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

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

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

	userRoute := r.Group("/users")
	{
		userRoute.POST("/", auth.Protect([]byte(secret), "STAFF"), userController.CreateUser)
		userRoute.GET("/", auth.Protect([]byte(secret), "STAFF"), userController.GetUsersByRole)
		userRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF", "TENANT"), userController.GetUserByID)
		userRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF", "TENANT"), userController.UpdateUser)
		userRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), userController.DeleteUser)
	}

	roomRepo := repository.NewRoomRepository(db)
	contractRepo := repository.NewContractRepository(db)
	roomService := service.NewRoomService(roomRepo, contractRepo)
	roomService.SetUserRepository(userRepo) // Initialize user repo for tenant operations
	roomController := controller.NewRoomController(roomService)

	roomRoute := r.Group("/rooms")
	{
		// CRUD Operations (STAFF only)
		roomRoute.POST("/", auth.Protect([]byte(secret), "STAFF"), roomController.CreateRoom)
		roomRoute.GET("/", auth.Protect([]byte(secret), "STAFF"), roomController.GetListRoom)
		roomRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomByID)
		roomRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.UpdateRoom)
		roomRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.DeleteRoom)

		// Relationship APIs (STAFF only)
		roomRoute.GET("/:id/contract", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomActiveContract)
		roomRoute.GET("/:id/contracts", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomContractHistory)
		roomRoute.GET("/:id/tenant", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomTenant)
		roomRoute.POST("/:id/assign", auth.Protect([]byte(secret), "STAFF"), roomController.AssignRoom)
	}

	contractService := service.NewContractService(contractRepo, roomRepo)
	contractService.SetUserRepository(userRepo)
	contractController := controller.NewContractController(contractService)

	contractRoute := r.Group("/contracts")
	{
		contractRoute.POST("/", auth.Protect([]byte(secret), "STAFF"), contractController.CreateContract)
		contractRoute.GET("/", auth.Protect([]byte(secret), "STAFF"), contractController.GetContracts)
		contractRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractByID)
		contractRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.UpdateContract)
		contractRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.DeleteContract)
		contractRoute.GET("/user/:userID", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractsByUserID)
		contractRoute.GET("/room/:roomID", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractsByRoomID)
	}

	meRoute := r.Group("/me")
	{
		// TENANT only endpoint
		meRoute.GET("/room", auth.Protect([]byte(secret), "TENANT"), roomController.GetMyRoom)
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
