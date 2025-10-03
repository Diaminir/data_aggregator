package server

import (
	"junior_effectivemobile/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server структура для наследования методов
type Server struct {
	app *handlers.HandlersApp
	log *logrus.Logger
}

// NewServer создает экземпляр структуры Server
func NewServer(app *handlers.HandlersApp, log *logrus.Logger) *Server {
	return &Server{
		app: app,
		log: log,
	}
}

// ServerStart настройка роутинга и запуск сервера
func (s *Server) ServerStart() *http.Server {
	router := gin.Default()
	api := router.Group("/subscriptions")
	{
		api.POST("", s.app.NewSubRecord)
		api.GET("", s.app.ListAllSubRecords)
		api.GET("/:id", s.app.GetUserSubRecord)
		api.PUT("/:id", s.app.UpdateSubRecord)
		api.DELETE("/:id", s.app.DeleteSubRecord)
		api.GET("/cost", s.app.CalculateCost)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка сервера")
		}
	}()
	return srv
}
