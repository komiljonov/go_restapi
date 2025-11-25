package api

import (
	"net/http"
	"restapi/db/sqlc"
	"restapi/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

type MeUpdateRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Birthdate   string `json:"birthdate"`
}

type AuthResponse struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birthdate"`

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) HandleRegister(c *gin.Context) {

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

	newUser, err := s.store.CreateUserTx(c.Request.Context(), db.CreateUserParams{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    passwordHash,
		Birthdate:   birthdate,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtAccess, jwtRefresh, err := newUser.CreateTokens()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := AuthResponse{
		Id:           newUser.ID,
		Name:         newUser.Name,
		PhoneNumber:  newUser.PhoneNumber,
		BirthDate:    newUser.Birthdate.Time.Format("2006-01-02"),
		AccessToken:  jwtAccess,
		RefreshToken: jwtRefresh,
	}

	c.JSON(http.StatusOK, resp)

}

func (s *Server) HandleLogin(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.store.GetByPhoneNumber(c.Request.Context(), req.PhoneNumber)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verified, err := user.CheckPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	jwtAccess, jwtRefresh, err := user.CreateTokens()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := AuthResponse{
		Id:           user.ID,
		Name:         user.Name,
		PhoneNumber:  user.PhoneNumber,
		BirthDate:    user.Birthdate.Time.Format("2006-01-02"),
		AccessToken:  jwtAccess,
		RefreshToken: jwtRefresh,
	}

	c.JSON(http.StatusOK, resp)

}

func (s *Server) HandleMe(c *gin.Context, user db.User) {

	c.JSON(http.StatusOK, user)

}

func (s *Server) HandleMeUpdate(c *gin.Context, user db.User) {
	var req MeUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := user

	if req.Name != "" {
		updated.Name = req.Name
	}

	if req.Birthdate != "" {
		bd, err := time.Parse("2006-01-02", req.Birthdate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid birthdate, expected format YYYY-MM-DD",
			})
			return
		}
		updated.Birthdate = pgtype.Date{Time: bd, Valid: true}

	}

	// 3) Persist to DB via your repository / store
	if _, err := s.store.UpdateUser(c.Request.Context(), db.UpdateUserParams{
		ID:        updated.ID,
		Name:      updated.Name,
		Birthdate: updated.Birthdate,
	}); err != nil {
		// adjust error handling as needed
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)

}
