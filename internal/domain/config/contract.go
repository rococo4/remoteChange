package config

import "remoteChange/internal/model"

type repo interface {
	CreateConfig(config model.ConfigEntity, userId int64) (int64, error)
	GetConfigByTeam(teamId int64) ([]model.ConfigEntity, error)
	UpdateConfig(config model.ConfigEntity, userId int64) (int64, error)
	Rollback(configId int64, userId int64) error
	GetConfigById(configId int64) (model.ConfigEntity, error)
	GetUserByUsername(username string) (model.UserEntity, error)
	GetIdOldConfig(configId int64) (int64, error)
}
type kuberClient interface {
	Deploy(entity model.ConfigEntity, toCreate bool) error
}
