package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("key")

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginHandler struct {
	usersFilePath string
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method: "+r.Method, http.StatusBadRequest)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	validPasswordHash, err := h.getPasswordHashByUser(user.Login)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(validPasswordHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Login: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

func (h *LoginHandler) getPasswordHashByUser(login string) (string, error) {
	usersData, err := os.ReadFile(h.usersFilePath)
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("unable to read users file")
	}
	var users []User
	err = json.Unmarshal(usersData, &users)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal users data")
	}

	for _, user := range users {
		if user.Login == login {
			return user.Password, nil
		}
	}
	return "", fmt.Errorf("user not found")
}
