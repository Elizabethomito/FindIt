package handlers

import (
	"encoding/json"
	"net/http"

	"findit/internal/db"
	"findit/internal/models"
	"findit/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

// Signup handles POST /signup
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if user.ID == "" || user.Email == "" || user.Password == "" {
		utils.JSONError(w, http.StatusBadRequest, "id, email, and password are required")
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (id, email, password) VALUES (?, ?, ?)",
		user.ID, user.Email, string(hashed),
	)
	if err != nil {
		utils.JSONError(w, http.StatusConflict, "user with this email or id already exists")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"message": "user created successfully",
	})
}

// Login handles POST /login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, email, password FROM users WHERE email = ?",
		creds.Email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "login successful",
		"user_id": user.ID,
	})
}
