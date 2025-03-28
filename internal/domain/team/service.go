package team

import (
	"context"
	"fmt"
	"remoteChange/internal/infrastructure"
	"remoteChange/internal/model"
)

type Service struct {
	repo teamRepo
}

func NewService(repo teamRepo) Service {
	return Service{repo: repo}
}

func (s *Service) CreateTeam(ctx context.Context, team model.TeamCreateDto) error {
	entity := infrastructure.MapTeamDtoCreateToTeamEntity(team)
	err := s.repo.SaveTeam(entity)
	if err != nil {
		return fmt.Errorf("error saving team: %w", err)
	}
	return nil
}

func (s *Service) EditUserInTeam(userId int64, teamId *int64) error {
	err := s.repo.UpdateUserInTeam(userId, teamId)
	if err != nil {
		return fmt.Errorf("error adding user to team: %w", err)
	}
	return nil
}

func (s *Service) GetTeamById(teamId int64) (*model.ResponseTeamDTO, error) {
	entity, err := s.repo.GetTeamById(teamId)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}
	dto := infrastructure.MapTeamEntityToResponseDto(entity)
	return &dto, nil
}

func (s *Service) EditTeam(team model.UpdateTeamDTO) error {
	entity, err := s.repo.GetTeamById(team.Id)
	if err != nil {
		return fmt.Errorf("error getting team: %w", err)
	}
	if team.Name != nil {
		entity.Name = *team.Name
	}
	if team.ClusterName != nil {
		entity.ClusterName = *team.ClusterName
	}
	if team.Namespace != nil {
		entity.Namespace = *team.Namespace
	}

	err = s.repo.EditTeam(entity)
	if err != nil {
		return fmt.Errorf("error editing team: %w", err)
	}
	return nil
}

func (s *Service) DeleteTeam(teamId int64) error {
	err := s.repo.DeleteTeam(teamId)
	if err != nil {
		return fmt.Errorf("error deleting team: %w", err)
	}
	return nil
}
