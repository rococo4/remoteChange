package team

import "remoteChange/internal/model"

type teamRepo interface {
	SaveTeam(team model.TeamEntity) (int64, error)
	UpdateUserInTeam(userId int64, teamId *int64) error
	GetTeamById(teamId int64) (model.TeamEntity, error)
	EditTeam(team model.TeamEntity) error
	DeleteTeam(teamId int64) error
	GetUserByUsername(username string) (model.UserEntity, error)
}
