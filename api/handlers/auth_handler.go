package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nishanth-thoughtclan/student-api/api/models"
	"github.com/nishanth-thoughtclan/student-api/api/services"
	"github.com/nishanth-thoughtclan/student-api/config"
	"github.com/nishanth-thoughtclan/student-api/utils"
)

func AuthHandler(cfg *config.Config, authService *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)

		if authService.ValidateUser(user) {
			token, err := utils.GenerateFirebaseToken(user.ID)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"token": token})
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
