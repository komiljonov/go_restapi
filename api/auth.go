package api

import (
	"context"
	"net/http"
	"strconv"

	"restapi/db/sqlc"
	"restapi/utils"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	Birthdate   string `json:"birthdate" binding:"required,datetime=2006-01-02"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birthdate"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) RegisterHandler(c *gin.Context) {

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birthdate, err := utils.ToPgDate(req.Birthdate)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordHash, err := utils.HashPasswordDefault(req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := s.store.CreateUserTx(context.Background(), db.CreateUserParams{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    passwordHash,
		Birthdate:   birthdate,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt_access, err_access := utils.GenerateJWT(strconv.Itoa(int(newUser.ID)), "access")
	jwt_refresh, err_refresh := utils.GenerateJWT(strconv.Itoa(int(newUser.ID)), "refresh")

	if err_access != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err_access.Error()})
		return
	}

	if err_refresh != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err_refresh.Error()})
		return
	}

	resp := RegisterResponse{
		Id:           newUser.ID,
		Name:         newUser.Name,
		PhoneNumber:  newUser.PhoneNumber,
		BirthDate:    newUser.Birthdate.Time.Format("2006-01-02"),
		AccessToken:  jwt_access,
		RefreshToken: jwt_refresh,
	}

	c.JSON(http.StatusOK, resp)

}

func (s *Server) HandleLogin(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.store.GetByPhoneNumber(context.Background())

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// verified := user.CheckPassword(req.Password)

	// if !verified {
	// }

}
