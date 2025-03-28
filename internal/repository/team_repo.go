package repository

import "remoteChange/internal/model"

func (r *Repo) SaveTeam(team model.TeamEntity) error {
	_, err := r.Db.Exec("insert into teams (name, cluster_name, namespace) values ($1, $2, $3)",
		team.Name, team.ClusterName, team.Namespace)
	return err
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
	_, err := r.Db.Exec("update teams set name=$1, cluster_name=$2, namespace=$3 where id=$4",
		team.Name, team.ClusterName, team.Namespace, team.Id)
	return err
}

func (r *Repo) DeleteTeam(teamId int64) error {
	_, err := r.Db.Exec("delete from teams where id=$1", teamId)
	return err
}
