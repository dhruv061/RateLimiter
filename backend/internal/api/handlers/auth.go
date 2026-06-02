package handlers

import (
	"fail2ban-dashboard/internal/auth"
	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	db  *database.DB
	jwt *auth.JWTManager
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *database.DB, jwt *auth.JWTManager) *AuthHandler {
	return &AuthHandler{db: db, jwt: jwt}
}

// Login authenticates a user and returns a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: username and password required")
		return
	}

	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, password_hash, email, role, must_change_pass FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.Role, &user.MustChangePass)
	if err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	token, err := h.jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.InternalError(c, "Failed to generate token")
		return
	}

	// Update last login
	h.db.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?", user.ID)

	response.OK(c, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetMe returns the current authenticated user.
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, email, role, must_change_pass, created_at, updated_at, last_login FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.MustChangePass, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.OK(c, user)
}

// ChangePassword updates the user's password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: current_password and new_password (min 6 chars) required")
		return
	}

	var currentHash string
	err := h.db.QueryRow("SELECT password_hash FROM users WHERE id = ?", userID).Scan(&currentHash)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		response.BadRequest(c, "Current password is incorrect")
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to hash password")
		return
	}

	_, err = h.db.Exec(
		"UPDATE users SET password_hash = ?, must_change_pass = 0, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		string(newHash), userID,
	)
	if err != nil {
		response.InternalError(c, "Failed to update password")
		return
	}

	response.Message(c, "Password updated successfully")
}

// Logout handles user logout (client-side token removal).
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Message(c, "Logged out successfully")
}
