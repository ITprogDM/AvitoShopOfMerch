package app

import (
	"ShopAvito/internal/config"
	"ShopAvito/internal/handlers"
	"ShopAvito/internal/repository"
	"ShopAvito/internal/services"
	"ShopAvito/pkg/logger"
	"ShopAvito/pkg/postgres"
	"ShopAvito/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	// Инициализация логгера
	log := logger.NewLogger()
	log.Info("Starting Avito Shop service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	db, err := postgres.ClientPostgres(log)
	if err != nil {

		log.Fatal("Failed to connect to database: ", err)
	}

	// Инициализация сервисов и обработчиков
	invenRepo := repository.NewInventoryRepository(db, log)
	userRepo := repository.NewUserRepository(db, log)
	invenService := services.NewInventoryService(invenRepo)
	authService := services.NewAuthService(userRepo, cfg.JwtSecret, log)
	userService := services.NewUserService(userRepo, log)
	transactionRepo := repository.NewTransactionRepository(db, log)
	transactionService := services.NewTransactionService(transactionRepo, userRepo, log)
	purchaseRepo := repository.NewPurchaseRepository(db, log)
	purchaseService := services.NewPurchaseService(purchaseRepo, userRepo, invenRepo, log)

	serv := new(server.Server)

	routes := handlers.RegisterRoutes(userService, transactionService, purchaseService, authService, invenService, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err = serv.RunServer(routes); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	log.Info("Сервер запущен на порту: 8080!")

	<-quit
	log.Info("Выключение сервера")

	if err = serv.ShutdownServer(); err != nil {
		log.Errorf("Ошибка при завершении работы сервера")
	}

	db.Close()
	log.Info("Сервер успешно выключен!")

}
