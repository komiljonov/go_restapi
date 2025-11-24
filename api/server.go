package api

import (
	"restapi/config"
	db "restapi/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config config.Config
	router *gin.Engine
	pool   *pgxpool.Pool
	store  db.Store
}

func NewServer(config config.Config, pool *pgxpool.Pool, store db.Store) *Server {
	engine := gin.Default()

	server := &Server{
		config: config,
		router: engine,
		pool:   pool,
		store:  store,
	}

	return server
}

func (s *Server) MountHandlers() {
	api := s.router.Group("/api")
	api.GET("/ping", s.Ping)

	auth := api.Group("/auth")
	auth.POST("/register", s.RegisterHandler)

	users := api.Group("/users")
	users.GET("", s.HandleAllUsers)

}

func (s *Server) Router() *gin.Engine {
	return s.router
}
