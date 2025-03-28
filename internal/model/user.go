package model

// TODO: добавить в дто валидацию

type UserDTORegister struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}

type UserDTOLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserEntity struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Surname  string `db:"surname"`
	TeamId   *int64 `db:"team_id"`
	Role     string `db:"role"`
	Password string `db:"password"`
}
