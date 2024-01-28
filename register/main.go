package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter login:")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hash := generateHashByPassword(password)
	saveUserToFile(login, hash)

}

func generateHashByPassword(pass string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(hash)
}

func saveUserToFile(login, password string) {
	usersData, err := os.ReadFile("/Users/egorskurihin/Desktop/Git/postgres-proxy/web-server/users.json")
	if err != nil {
		panic(err)
	}
	users := []User{}
	err = json.Unmarshal(usersData, &users)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	users = append([]User{{Login: login, Password: password}}, users...)
	jsonUsers, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("/Users/egorskurihin/Desktop/Git/postgres-proxy/web-server/users.json", jsonUsers, 0644)
	if err != nil {
		panic(err)
	}
	log.Println("User successfuly created")
}
