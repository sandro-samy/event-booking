package models

import (
	"errors"

	"github.com/sandro-samy/event-booking/db"
	"github.com/sandro-samy/event-booking/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() error {
	query := `INSERT INTO users(email, password) VALUES (?, ?)`

	stmt, err := db.DB.Prepare(query)

	defer stmt.Close()
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(user.Password)


	result, err := stmt.Exec(user.Email, hashedPassword)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = userID

	return nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var password string
	err := row.Scan(&u.ID, &password)

	if err != nil {
		return errors.New("credentials invalid")
	}

	isValid := utils.CheckMatchingPassword(password, u.Password)

	if !isValid {
		return errors.New("credentials invalid")
	}

	return nil
}
