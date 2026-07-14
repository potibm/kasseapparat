package models

type AuthUser struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
