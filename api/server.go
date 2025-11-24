package api

import (
	"restapi/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config config.Config
	router *gin.Engine
	pool   *pgxpool.Pool
}

func NewServer(config config.Config, pool *pgxpool.Pool) *Server {
	engine := gin.Default()

	server := &Server{
		config: config,
		router: engine,
		pool:   pool,
	}

	return server
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")
	api.GET("/ping", s.Ping)

	auth := api.Group("/auth")
	auth.POST("/register", s.RegisterHandler)

}

func (s *Server) Router() *gin.Engine {
	return s.router
}
