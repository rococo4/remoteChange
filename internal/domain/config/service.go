package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	auth "remoteChange/internal/domain/jwt"
	"remoteChange/internal/infrastructure"
	"remoteChange/internal/model"
)

type Service struct {
	configRepo  repo
	kuberClient kuberClient
}

func NewService(configRepo repo, client kuberClient) *Service {
	return &Service{configRepo: configRepo, kuberClient: client}
}

func (s *Service) CreateConfig(ctx context.Context, config model.CreateConfigDTO) error {
	user, err := s.getUserFromCtx(ctx)
	if err != nil || user == nil {
		return fmt.Errorf("error getting user from ctx %w", err)
	}
	configEntity := infrastructure.MapCreateConfigDtoToConfigEntity(config)
	if user.TeamId == nil {
		return fmt.Errorf("user has no team")
	}
	configEntity.TeamId = *user.TeamId
	_, err = s.configRepo.CreateConfig(configEntity, user.Id)
	//err = s.kuberClient.Deploy(configEntity, true)
	if err != nil {
		return fmt.Errorf("error deploying config %w", err)
	}

	return err
}

func (s *Service) EditConfig(ctx context.Context, config model.UpdateConfigDTO) error {
	actualConfigId := config.Id
	configId, err := s.configRepo.GetIdByConfigVersionId(config.Id)
	if err != nil {
		return fmt.Errorf("error getting actual config id by config id %w", err)
	}
	config.Id = configId

	user, err := s.getUserFromCtx(ctx)
	if err != nil || user == nil {
		return fmt.Errorf("error getting user from ctx %w", err)
	}
	configEntity, err := s.configRepo.GetConfigById(config.Id)
	if err != nil {
		return fmt.Errorf("error getting config by id %w", err)
	}
	configEntity.Content = config.Content
	//err = s.kuberClient.Deploy(configEntity, false)
	if err != nil {
		return fmt.Errorf("error deploying config %w", err)
	}
	_, err = s.configRepo.UpdateConfig(configEntity, user.Id, actualConfigId)
	if err != nil {
		return fmt.Errorf("error updating config %w", err)
	}
	return nil
}

func (s *Service) GetConfigByTeam(teamId int64) ([]model.ConfigResponse, error) {
	configs, err := s.configRepo.GetConfigByTeam(teamId)
	if err != nil {
		return nil, fmt.Errorf("error getting configs by team %w", err)
	}
	configsResponse := make([]model.ConfigResponse, 0, len(configs))
	for _, config := range configs {
		normalId, err := s.configRepo.GetActualConfigIdByConfigId(config.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error getting actual config id by config id %w", err)
		}
		if !errors.Is(err, sql.ErrNoRows) {
			config.Id = normalId
			configsResponse = append(configsResponse, infrastructure.MapConfigEntityToConfigResponse(config))
		}
	}
	return configsResponse, nil
}

func (s *Service) Rollback(ctx context.Context, configId int64) error {
	user, err := s.getUserFromCtx(ctx)
	if err != nil || user == nil {
		return fmt.Errorf("error getting user from ctx %w", err)
	}

	configId, err = s.configRepo.GetIdByConfigVersionId(configId)
	if err != nil {
		return fmt.Errorf("error getting actual config id by config id %w", err)
	}

	configIdToDeploy, err := s.configRepo.GetIdOldConfig(configId)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error getting old config id %w", err)
	}
	if configIdToDeploy == nil {
		return fmt.Errorf("nothing to rollback")
	}
	configEntity, err := s.configRepo.GetConfigById(*configIdToDeploy)
	if err != nil {
		return fmt.Errorf("error getting config by id %w", err)
	}
	configEntity.Description = "time.Now()"
	//err = s.kuberClient.Deploy(configEntity, false)
	if err != nil {
		return fmt.Errorf("error deploying config %w", err)
	}

	return s.configRepo.Rollback(configId, user.Id)
}

func (s *Service) GetConfigChangesForId(configId int64) ([]model.ConfigVersionResponse, error) {
	id, err := s.configRepo.GetIdByConfigVersionId(configId)
	if err != nil {
		return nil, fmt.Errorf("error getting id by config version id %w", err)
	}
	var configChanges []model.ConfigChangesEntity
	oldConfigId, err := s.configRepo.GetOldConfigId(id)
	if err != nil {
		return nil, err
	}
	for {
		if oldConfigId == nil {
			break
		}
		configChange, err := s.configRepo.GetConfigChanges(*oldConfigId)
		if err != nil {
			return nil, fmt.Errorf("error getting config changes %w", err)
		}
		oldConfigId = configChange.OldConfig
		configChanges = append(configChanges, configChange)
	}
	var resp []model.ConfigVersionResponse
	for _, change := range configChanges {
		res, _ := infrastructure.MapConfigChangesEntityToConfigVersionResponse(change, s.configRepo)
		resp = append(resp, res)
	}
	return resp, nil
}

func (s *Service) GetActualConfigById(configId int64) (model.ConfigResponse, error) {
	actualConfigId, err := s.configRepo.GetIdByConfigVersionId(configId)
	if err != nil {
		return model.ConfigResponse{}, fmt.Errorf("error getting actual config id by config id %w", err)
	}
	config, err := s.configRepo.GetConfigById(actualConfigId)
	if err != nil {
		return model.ConfigResponse{}, fmt.Errorf("error getting config by id %w", err)
	}
	return infrastructure.MapConfigEntityToConfigResponse(config), nil
}

func (s *Service) getUserFromCtx(ctx context.Context) (*model.UserEntity, error) {
	claims, ok := ctx.Value("user").(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("error getting claims")
	}
	user, err := s.configRepo.GetUserByUsername(claims.Username)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username %w", err)
	}
	return &user, nil
}
