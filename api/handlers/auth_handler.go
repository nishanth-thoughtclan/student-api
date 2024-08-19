package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/utils"
)

type UserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

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
func AuthHandler(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest UserRequest
		// Decode the incoming user data
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err := validate.Struct(userRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Retrieve the user's ID
		userID, err := authService.GetUserIDByEmail(r.Context(), userRequest.Email)
		if err != nil {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}

		// Validate the user credentials
		validUser, err := authService.ValidateUser(r.Context(), userRequest.Email, userRequest.Password)
		if err != nil || !validUser {
			http.Error(w, "Invalid email/password", http.StatusUnauthorized)
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

func SingUpHandler(authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest UserRequest
		// Decode the incoming user data
		if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err := validate.Struct(userRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := userFromAddUserRequest(userRequest)

		updatedUser, err := authService.CreateUser(r.Context(), user)
		if err != nil {
			http.Error(w, "User Already Exists with the email", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(userToUserResponse(*updatedUser)); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func userFromAddUserRequest(r UserRequest) models.User {
	return models.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

func userToUserResponse(u models.User) UserResponse {
	return UserResponse{
		Email: u.Email,
		Id:    u.ID.String(),
	}
}
