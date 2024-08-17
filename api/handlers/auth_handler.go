package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/config"
	"github.com/nishanth-thoughtclan/student-api/utils"
)

// AuthHandler handles user authentication
// @Summary Authenticate user
// @Description Authenticates a user and returns a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param authRequest body AuthRequest true "User credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth [post]
func AuthHandler(cfg *config.Config, authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		// Decode the incoming user data
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the user credentials
		validUser, err := authService.ValidateUser(user.Email, user.Password)
		if err != nil || !validUser {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve the user's ID
		userID, err := authService.GetUserIDByEmail(user.Email)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Generate a JWT token with the user ID as the subject
		token, err := utils.GenerateToken(userID.String())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Return the token to the client
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
