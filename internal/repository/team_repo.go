package repository

import "remoteChange/internal/model"

func (r *Repo) SaveTeam(team model.TeamEntity) (int64, error) {
	var id int64
	err := r.Db.QueryRow("insert into teams (name) values ($1) returning id", team.Name).Scan(&id)
	return id, err
}

func (r *Repo) UpdateUserInTeam(userId int64, teamId *int64) error {
	if teamId == nil {
		_, err := r.Db.Exec("update users set team_id = null where id = $1", userId)
		return err
	}
	_, err := r.Db.Exec("update users set team_id=$1 where id=$2", teamId, userId)
	return err
}

func (r *Repo) GetTeamById(teamId int64) (model.TeamEntity, error) {
	var team model.TeamEntity
	err := r.Db.Get(&team, "select * from teams where id=$1", teamId)
	return team, err
}

func (r *Repo) EditTeam(team model.TeamEntity) error {
	_, err := r.Db.Exec("update teams set name=$1 where id=$2",
		team.Name, team.Id)
	return err
}

func (r *Repo) DeleteTeam(teamId int64) error {
	_, err := r.Db.Exec("delete from teams where id=$1", teamId)
	return err
}
