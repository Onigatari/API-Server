package handler

import (
	_ "Avito/docs"
	"Avito/internal/repository"
	"Avito/internal/service"
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

type Handler struct {
	services *service.BillingService
}

func NewHandler(services *service.BillingService) *Handler {
	return &Handler{services: services}
}

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(Port string) error {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatalf("[Handler] Can't connect to database: %s", err.Error())
	} else {
		log.Println("Database connection!")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	billingHandler := NewHandler(services)
	billingRouter := billingHandler.configureRoutes()

	s.httpServer = &http.Server{
		Addr:    ":" + Port,
		Handler: billingRouter,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (h *Handler) configureRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/ping", h.ping)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	accountChanges := router.Group("/api/")
	{
		accountChanges.GET("/balance/:id", h.balance)
		accountChanges.POST("/deposit", h.deposit)
		accountChanges.POST("/withdrawal", h.withdrawal)
		accountChanges.POST("/transfer", h.transfer)

		accountChanges.POST("/reserveService", h.reserveService)
		accountChanges.POST("/approveService", h.approveOrder)
		accountChanges.POST("/failedService", h.failedService)
	}

	return router
}
