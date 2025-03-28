package repository

import (
	"remoteChange/internal/model"
)

func (r *Repo) SaveUser(user model.UserEntity) error {
	_, err := r.Db.Exec("insert into users (username, name, surname, role, password) values ($1, $2, $3, $4, $5)",
		user.Username, user.Name, user.Surname, user.Role, user.Password)
	return err
}
func (r *Repo) GetUserByUsername(username string) (model.UserEntity, error) {
	var user model.UserEntity
	err := r.Db.Get(&user, "select * from users where username=$1", username)
	return user, err
}
