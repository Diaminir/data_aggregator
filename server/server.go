package server

import (
	"junior_effectivemobile/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	app *handlers.HandlersApp
}

func NewServer(app *handlers.HandlersApp) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) ServerStart() {
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Ошибка сервера")
	}
}
