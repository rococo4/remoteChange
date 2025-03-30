package model

type TeamCreateDto struct {
	Name string `json:"name"`
}

type TeamEntity struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type UpdateTeamDTO struct {
	Id   int64   `json:"id"`
	Name *string `json:"name"`
}
type ResponseTeamDTO struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
