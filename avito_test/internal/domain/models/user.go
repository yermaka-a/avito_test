package models

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func NewUser(id, username, TeamName string, isActive bool) *User {
	return &User{
		UserID:   id,
		Username: username,
		TeamName: TeamName,
		IsActive: isActive,
	}
}
