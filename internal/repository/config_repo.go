package repository

import (
	"remoteChange/internal/model"
	"time"
)

func (r *Repo) CreateConfig(config model.ConfigEntity, userId int64) (int64, error) {
	var id int64
	err := r.Db.QueryRow("insert into configs (team_id, type, content, created_at) values ($1, $2, $3, $4) returning id",
		config.TeamId, config.Type, config.Content, config.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	_, err = r.Db.Exec("insert into config_versions (actual_config_id) values ($1)", id)
	if err != nil {
		return 0, err
	}
	_, err = r.Db.Exec("insert into config_changes (new_config, old_config, user_id, action, action_at, team_id) values ($1, $2, $3, $4, $5, $6)",
		id, nil, userId, "create", time.Now(), config.TeamId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) GetConfigByTeam(teamId int64) ([]model.ConfigEntity, error) {
	var configs []model.ConfigEntity
	err := r.Db.Select(&configs, "select * from configs where team_id=$1", teamId)
	return configs, err
}

func (r *Repo) GetConfigById(configId int64) (model.ConfigEntity, error) {
	var config model.ConfigEntity
	err := r.Db.Get(&config, "select * from configs where id=$1", configId)
	return config, err
}

func (r *Repo) SaveConfig(config model.ConfigEntity) error {
	_, err := r.Db.Exec("update configs set team_id=$1, type=$2, content=$3, created_at=$4, name=$5 where id=$6",
		config.TeamId, config.Type, config.Content, config.CreatedAt, config.Name, config.Id)
	return err
}

func (r *Repo) UpdateConfig(config model.ConfigEntity, userId int64, configVersionId int64) (int64, error) {
	var newActualConfigId int64
	// вставляем в таблицу configs новый конфиг
	err := r.Db.QueryRow("insert into configs (team_id, type, content, created_at, name) values ($1, $2, $3, $4, $5) returning id",
		config.TeamId, config.Type, config.Content, config.CreatedAt, config.Name).Scan(&newActualConfigId)
	if err != nil {
		return 0, err
	}
	// получаем id(в таблице configs) старого конфига из таблицы config_versions по id
	var oldActualConfigId int64
	err = r.Db.QueryRow("select actual_config_id from config_versions where id=$1", configVersionId).Scan(&oldActualConfigId)
	if err != nil {
		return 0, err
	}
	// обновляем actual_config_id в таблице config_versions на новый конфиг id
	_, err = r.Db.Exec("update config_versions set actual_config_id=$1 where id=$2", newActualConfigId, configVersionId)
	if err != nil {
		return 0, err
	}
	var oldId int64
	// получаем id старого конфига из таблицы config_changes
	err = r.Db.QueryRow("select id from config_changes where new_config=$1", oldActualConfigId).Scan(&oldId)
	if err != nil {
		return 0, err
	}
	// добавляем запись в таблицу config_changes об изменении конфига
	_, err = r.Db.Exec("insert into config_changes (new_config, old_config, user_id, action, action_at, team_id) values ($1, $2, $3, $4, $5, $6)",
		newActualConfigId, oldId, userId, "update", time.Now(), config.TeamId)
	if err != nil {
		return 0, err
	}
	return newActualConfigId, nil
}

func (r *Repo) GetIdOldConfig(configId int64) (*int64, error) {
	var oldId *int64
	err := r.Db.QueryRow("select old_config from config_changes where new_config=$1", configId).Scan(&oldId)
	if err != nil {
		return nil, err
	}
	var oldConfigId *int64
	err = r.Db.QueryRow("select new_config from config_changes where id=$1", oldId).Scan(&oldConfigId)
	if err != nil {
		return nil, err
	}
	return oldConfigId, nil
}

func (r *Repo) Rollback(configId int64, userId int64) error {
	var oldId, teamId int64
	// получаем id старого конфига для того, который нужно rollback. Id из таблицы config_changes
	err := r.Db.QueryRow("select old_config, team_id from config_changes where new_config=$1", configId).Scan(&oldId, &teamId)
	if err != nil {
		return err
	}
	var oldConfigId int64
	// получаем Id старого конфига из таблицы configs
	err = r.Db.QueryRow("select new_config from config_changes where id=$1", oldId).Scan(&oldConfigId)
	if err != nil {
		return err
	}
	// обновляем actual_config_id в таблице config_versions
	_, err = r.Db.Exec("update config_versions set actual_config_id=$1 where actual_config_id=$2", oldConfigId, configId)
	if err != nil {
		return err
	}
	// добавляем запись в таблицу config_changes об изменении конфига
	_, err = r.Db.Exec("insert into config_changes (new_config, old_config, user_id, action, action_at, team_id) values ($1, $2, $3, $4, $5, $6)",
		oldConfigId, configId, userId, "rollback", time.Now(), teamId)
	return err
}

// получение актуального id в таблице configs по актуальному id конфига
func (r *Repo) GetIdByConfigVersionId(configVersionId int64) (int64, error) {
	var id int64
	err := r.Db.QueryRow("select actual_config_id from config_versions where id=$1", configVersionId).Scan(&id)
	return id, err
}

func (r *Repo) GetOldConfigId(newConfigId int64) (*int64, error) {
	var oldConfigId *int64
	err := r.Db.Get(&oldConfigId, "select old_config from config_changes where new_config=$1", newConfigId)
	return oldConfigId, err
}

func (r *Repo) GetConfigChanges(configId int64) (model.ConfigChangesEntity, error) {
	var configVersions model.ConfigChangesEntity
	err := r.Db.Get(&configVersions, "select * from config_changes where id=$1", configId)
	return configVersions, err
}

// получение актульного id в таблице config_versions по actual_config_id конфига
func (r *Repo) GetActualConfigIdByConfigId(id int64) (int64, error) {
	var actualId int64
	err := r.Db.Get(&actualId, "select id from config_versions where actual_config_id=$1", id)
	return actualId, err
}
