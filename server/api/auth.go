package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	m "social-network/models"
	"social-network/pkg/db/sqlite"
	"social-network/util"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user m.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error reading data", 400)
		return
	}

	// check if all the required fields are provided
	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Password) == "" || strings.TrimSpace(user.FirstName) == "" || strings.TrimSpace(user.LastName) == "" || user.DateOfBirth.IsZero() {
		http.Error(w, "Please populate all required fields", 400)
		return
	}

	var id int

	// check if the username or email already exists
	err := sqlite.DB.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", user.Email, user.Username).Scan(&id)

	if err == nil {
		http.Error(w, "User already exists", 400)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "Something went wrong", 500)
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Something went wrong", 400)
		return
	}

	if _, err := sqlite.DB.Exec("INSERT INTO users (username, email, password, firstname, lastname, date_of_birth) VALUES (?, ?, ?, ?, ?, ?)", user.Username, user.Email, string(hashedpassword), user.FirstName, user.LastName, user.DateOfBirth); err != nil {
		http.Error(w, "Something went wrong, please try again later", 500)
		return
	}

	// generate the session for the user
	if err := util.GenerateSession(w, &user); err != nil {
		http.Error(w, "Something went wrong", 500)
		log.Printf("Error: %v", err)
		return
	}

	w.Write([]byte("User regsitered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("user login")
}
