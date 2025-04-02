package model

import "time"

type ConfigEntity struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	TeamId      int64     `db:"team_id"`
	Type        string    `db:"type"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	Description string    `db:"description"`
}

type CreateConfigDTO struct {
	TeamId      int64  `json:"team_id"`
	Content     string `json:"content"`
	Description string `json:"description"`
}
type UpdateConfigDTO struct {
	Id      int64  `db:"id"`
	Content string `db:"content"`
}

type ConfigResponse struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	TeamId      int64     `json:"team_id"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}

type ConfigVersionResponse struct {
	Id       int64     `json:"id"`
	Username string    `json:"username"`
	Action   string    `json:"action"`
	Date     time.Time `json:"date"`
}

type ConfigChangesEntity struct {
	Id        int64     `db:"id"`
	NewConfig int64     `db:"new_config"`
	OldConfig *int64    `db:"old_config"`
	UserId    int64     `db:"user_id"`
	TeamId    int64     `db:"team_id"`
	Action    string    `db:"action"`
	ActionAt  time.Time `db:"action_at"`
}
