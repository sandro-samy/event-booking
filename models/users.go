package models

import (
	"github.com/google/uuid"

	"github.com/sandro-samy/event-booking/utils"
	"github.com/sandro-samy/event-booking/db"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) Save() error {
	query := `INSERT INTO users(id, email, password) VALUES (?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	
	if err != nil {
		return err
	}

	id := uuid.NewString()

	_, err = stmt.Exec(id, user.Email, hashedPassword)
	if err != nil {
		return err
	}

	user.ID = id

	return nil
}
