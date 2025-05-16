package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` //Insecure! Just for test
	CreatedAt time.Time `json:"created_at"`
}
