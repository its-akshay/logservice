package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/logservice/internal/model"
	"github.com/logservice/internal/repo"
	"github.com/logservice/internal/utils"
)

type AuthHandler struct {
	repo repo.LogRepository
}

func NewAuthHandler(r repo.LogRepository) *AuthHandler {
	return &AuthHandler{repo: r}
}

// Register godoc
// @Summary Register user
// @Description Create new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Register body"
// @Success 201 {string} string "created"
// @Failure 400 {object} model.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  "invalid_request",
		})
		return
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, model.ErrorResponse{
			Error: err.Error(),
			Code:  "internal_error",
		})
		return
	}
	user := model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Role:         req.Role,
	}

	if err := h.repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: err.Error(),
			Code:  "internal_error",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login body"
// @Success 200 {object} map[string]string
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "invalid body",
			Code:  "INVALID_BODY",
		})
		return
	}

	user, err := h.repo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: "unauthorized",
			Code:  "UNAUTHORIZED",
		})
		return
	}

	isValid := utils.CheckPassword(user.PasswordHash, req.Password)
	if !isValid {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: "unauthorized",
			Code:  "UNAUTHORIZED",
		})
		return
	}

	// Generate JWT tokens
	accessToken, err := repo.GenerateAccessToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "internal error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	refreshToken, err := repo.GenerateRefreshToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "internal error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RefreshRequest true "Refresh body"
// @Success 200 {object} map[string]string
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "invalid body",
			Code:  "INVALID_BODY",
		})
		return
	}

	claims, err := repo.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: "invalid refresh token",
			Code:  "INVALID_TOKEN",
		})
		return
	}
	userID, _ := uuid.Parse(claims.Subject)

	user := model.User{
		ID:       userID,
		Username: claims.Username,
		Role:     claims.Role,
	}
	accessToken, err := repo.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "internal error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}

	refreshToken, err := repo.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: "internal error",
			Code:  "INTERNAL_ERROR",
		})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}
