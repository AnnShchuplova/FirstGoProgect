package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"FurryTrack/internal/config"
	"FurryTrack/internal/controllers"
	"FurryTrack/internal/repositories"
	"FurryTrack/internal/services"
	"FurryTrack/pkg/database"
	"FurryTrack/pkg/middleware"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к БД через GORM
	gormDB, err := database.ConnectGORM(database.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Настройка соединения
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to database")

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(gormDB)
	petRepo := repositories.NewPetRepository(gormDB)
	vaccineRepo := repositories.NewVaccineRepository(gormDB)
	vaccineRecordRepo := repositories.NewVaccineRecordRepository(gormDB)
	postRepo := repositories.NewPostRepository(gormDB)
	userRelationRepo := repositories.NewUserRelationRepository(gormDB)
	eventRepo := repositories.NewEventRepository(gormDB)


	// Инициализация сервисов
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret)
	petService := services.NewPetService(petRepo)
	vaccineService := services.NewVaccineService(vaccineRepo)

	vaccineRecordService := services.NewVaccineRecordService(
		vaccineRecordRepo,
		petRepo,    
		vaccineRepo, 
	)

	postService := services.NewPostService(postRepo, petRepo)
	userRelationService := services.NewUserRelationService(userRelationRepo)
	feedService := services.NewFeedService(postRepo, userRelationRepo)
	eventService := services.NewEventService(eventRepo)

	// Инициализация контроллеров
	authController := controllers.NewAuthController(authService)
	petController := controllers.NewPetController(*petService, *vaccineService)
	vaccineController := controllers.NewVaccineController(*vaccineService)
	vaccineRecordController := controllers.NewVaccineRecordController(*vaccineRecordService)
	postController := controllers.NewPostController(*postService)
	userRelationController := controllers.NewUserRelationController(userRelationService)
	feedController := controllers.NewFeedController(feedService, userRelationService)
	eventController := controllers.NewEventController(eventService)

	

	// Настройка Gin
	router := gin.Default()

	// Настройка CORS
	router.Use(middleware.GinCorsMiddleware())

	// Публичные маршруты (не требуют аутентификации)
	public := router.Group("/api")
	{
		// Аутентификация
		public.POST("/register", authController.RegisterUser)
		public.POST("/login", authController.LoginUser)
	}

	// Приватные маршруты (требуют JWT)
	private := router.Group("/api")
	private.Use(middleware.GinAuthMiddleware(cfg.JWT.Secret))
	{
		// Профиль пользователя
		private.GET("/profile", authController.GetUserProfile)
		private.GET("/users/email/:email", authController.GetUserByEmail)

		// Питомцы
		private.POST("/pets", petController.CreatePet)
		private.GET("/pets/:pet_id", petController.GetPet)
		private.GET("/pets", petController.GetUserPets)
		private.PUT("/pets/:pet_id", petController.UpdatePet) 
		private.DELETE("/pets/:pet_id", petController.DeletePet)
		private.POST("/pets/:pet_id/photo", petController.UploadPetPhoto)

		// Вакцины как сущности 
		private.GET("/vaccines", vaccineController.GetAllVaccines)
		private.POST("/vaccines", vaccineController.CreateVaccine)

		// Записи о вакцинации
		private.POST("/pets/:pet_id/vaccine-records", vaccineRecordController.AddVaccineRecord)
		private.GET("/pets/:pet_id/vaccine-records", vaccineRecordController.GetPetVaccineHistory)

		// Посты
		private.POST("/posts", postController.CreatePost)
		private.GET("/posts/feed", postController.GetFeed)

		// Подписки (User Relations)
        private.POST("/users/:id/follow", userRelationController.Follow)
        private.GET("/users/me/following", userRelationController.GetFollowing)
        private.GET("/users/me/followers", userRelationController.GetFollowers)
        private.GET("/users/:id/following", userRelationController.GetFollowing)
        private.GET("/users/:id/followers", userRelationController.GetFollowers)

		// Лента 
		private.GET("/feed/main", feedController.GetMainFeed)
		private.GET("/feed/market", feedController.GetMarketFeed)
		private.GET("/feed/following", feedController.GetFollowingFeed)
		
		// Лайки и комментарии 
		private.POST("/posts/:id/like", postController.LikePost)
		private.POST("/posts/:id/comments", postController.AddComment)
		private.GET("/posts/:id/comments", postController.GetComments)

		// События 
		private.POST("/events", eventController.CreateEvent)
        private.GET("/pets/:pet_id/events", eventController.GetPetEvents)
	}

	// Настройка сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	log.Printf("Server started on port %s", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Таймаут для завершения операций
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
