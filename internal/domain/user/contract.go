package user

import "remoteChange/internal/model"

type userRepo interface {
	SaveUser(user model.UserEntity) error
	GetUserByUsername(username string) (model.UserEntity, error)
}
