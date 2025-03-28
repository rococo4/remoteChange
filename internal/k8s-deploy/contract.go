package k8s_deploy

import "remoteChange/internal/model"

type repo interface {
	GetTeamById(teamId int64) (model.TeamEntity, error)
	SaveConfig(config model.ConfigEntity) error
}
