package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PunMung-66/ApartmentSys/config"
	"github.com/PunMung-66/ApartmentSys/controller"
	"github.com/PunMung-66/ApartmentSys/internal/auth"
	"github.com/PunMung-66/ApartmentSys/internal/storage"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/repository"
	"github.com/PunMung-66/ApartmentSys/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	env := flag.String("e", "", "environment (local|dev)")
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

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"https://apartment-managment-front-end.vercel.app",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{
			"status":  http.StatusOK,
			"version": "0.0.1",
		})
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
			&model.BillSlip{},
		)

		fmt.Println("Setup completed!")
		return
	}

	// ================= USER =================
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userRoute := r.Group("/users")
	{
		userRoute.POST("", auth.Protect([]byte(secret), "STAFF"), userController.CreateUser)
		userRoute.GET("", auth.Protect([]byte(secret), "STAFF"), userController.GetUsersByRole)
		userRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF", "TENANT"), userController.GetUserByID)
		userRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF", "TENANT"), userController.UpdateUser)
		userRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), userController.DeleteUser)
	}

	// ================= ROOM =================
	roomRepo := repository.NewRoomRepository(db)
	contractRepo := repository.NewContractRepository(db)

	roomService := service.NewRoomService(roomRepo, contractRepo)
	roomService.SetUserRepository(userRepo)

	roomController := controller.NewRoomController(roomService)

	roomRoute := r.Group("/rooms")
	{
		roomRoute.POST("", auth.Protect([]byte(secret), "STAFF"), roomController.CreateRoom)
		roomRoute.GET("", auth.Protect([]byte(secret), "STAFF"), roomController.GetListRoom)
		roomRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomByID)
		roomRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.UpdateRoom)
		roomRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), roomController.DeleteRoom)

		roomRoute.GET("/:id/contract", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomActiveContract)
		roomRoute.GET("/:id/contracts", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomContractHistory)
		roomRoute.GET("/:id/tenant", auth.Protect([]byte(secret), "STAFF"), roomController.GetRoomTenant)
		roomRoute.POST("/:id/assign", auth.Protect([]byte(secret), "STAFF"), roomController.AssignRoom)
	}

	// ================= CONTRACT =================
	contractService := service.NewContractService(contractRepo, roomRepo)
	contractService.SetUserRepository(userRepo)

	contractController := controller.NewContractController(contractService)

	contractRoute := r.Group("/contracts")
	{
		contractRoute.POST("", auth.Protect([]byte(secret), "STAFF"), contractController.CreateContract)
		contractRoute.GET("", auth.Protect([]byte(secret), "STAFF"), contractController.GetContracts)
		contractRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractByID)
		contractRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.UpdateContract)
		contractRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), contractController.DeleteContract)
		contractRoute.GET("/user/:userID", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractsByUserID)
		contractRoute.GET("/room/:roomID", auth.Protect([]byte(secret), "STAFF"), contractController.GetContractsByRoomID)
	}

	// ================= UTILITY RATE =================
	utilityRateRepo := repository.NewUtilityRateRepository(db)
	utilityRepo := repository.NewUtilityUsageRepository(db)
	utilityService := service.NewUtilityService(utilityRateRepo, utilityRepo, contractRepo)
	utilityController := controller.NewUtilityController(utilityService)

	utilityRateRoute := r.Group("/utility-rates")
	{
		utilityRateRoute.POST("", auth.Protect([]byte(secret), "STAFF"), utilityController.CreateRate)
		utilityRateRoute.GET("", auth.Protect([]byte(secret), "STAFF", "TENANT"), utilityController.GetAllRates)
		utilityRateRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF", "TENANT"), utilityController.GetRateByID)
		utilityRateRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), utilityController.UpdateRate)
		utilityRateRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), utilityController.DeleteRate)
	}

	utilityUsageRoute := r.Group("/utility-usages")
	{
		utilityUsageRoute.POST("", auth.Protect([]byte(secret), "STAFF"), utilityController.RecordUsage)
		utilityUsageRoute.GET("/contract/:contractID", auth.Protect([]byte(secret), "STAFF"), utilityController.GetUsagesByContract)
		utilityUsageRoute.PUT("/:id", auth.Protect([]byte(secret), "STAFF"), utilityController.UpdateUsage)
		utilityUsageRoute.DELETE("/:id", auth.Protect([]byte(secret), "STAFF"), utilityController.DeleteUsage)
		utilityUsageRoute.GET("/:id", auth.Protect([]byte(secret), "STAFF"), utilityController.GetUsageByID)
	}

	// ================= BILL SLIP (Supabase) =================
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")
	if supabaseKey == "" {
		fmt.Println("ERROR: Supabase Key is completely empty! Check .env.dev")
	} else {
		fmt.Println("Supabase Key loaded successfully!")
	}
	storageClient := storage.NewSupabaseStorage(supabaseURL, supabaseKey)

	billSlipRepo := repository.NewBillSlipRepository(db)
	billSlipService := service.NewBillSlipService(billSlipRepo, storageClient)
	billSlipController := controller.NewBillSlipController(billSlipService)

	billSlipRoute := r.Group("/billslips")
	{
		billSlipRoute.POST("/upload", billSlipController.UploadBillSlip)
	}

	// ================= ME =================
	meRoute := r.Group("/me")
	{
		meRoute.GET("/room", auth.Protect([]byte(secret), "TENANT"), roomController.GetMyRoom)
		meRoute.GET("/usages", auth.Protect([]byte(secret), "TENANT"), utilityController.GetMyUsages)
		meRoute.GET("/usages/latest", auth.Protect([]byte(secret), "TENANT"), utilityController.GetMyLatestUsage)
	}

	// ================= AUTH =================
	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService, []byte(secret))

	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", authController.LoginHandler)
		authRoute.POST("/register", authController.RegisterHandler)
	}

	r.Run(":" + port)
}