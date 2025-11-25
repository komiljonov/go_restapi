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

	auth := api.Group("/auth", s.AuthMiddleware())
	auth.POST("/register", s.HandleRegister)
	auth.POST("/login", s.HandleLogin)
	auth.GET("/me", s.AuthRequired(s.HandleMe))
	auth.PATCH("/me", s.AuthRequired(s.HandleMeUpdate))

	users := api.Group("/users", s.AuthMiddleware())
	users.GET("", s.HandleAllUsers)

}

func (s *Server) Router() *gin.Engine {
	return s.router
}
