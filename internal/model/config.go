package model

import "time"

type ConfigEntity struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	TeamId    int64     `db:"team_id"`
	Type      string    `db:"type"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateConfigDTO struct {
	TeamId    int64     `json:"team_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
type UpdateConfigDTO struct {
	Id      int64  `db:"id"`
	Content string `db:"content"`
}

type ConfigResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TeamId    int64     `json:"team_id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
