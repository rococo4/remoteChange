package team

import (
	"context"
	"fmt"
	auth "remoteChange/internal/domain/jwt"
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
	user, err := s.getUserFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("error getting user from ctx: %w", err)
	}
	id, err := s.repo.SaveTeam(entity)
	if err != nil {
		return fmt.Errorf("error saving team: %w", err)
	}
	user.TeamId = &id
	err = s.repo.UpdateUserInTeam(user.Id, user.TeamId)
	if err != nil {
		fmt.Errorf("error updating user in team: %w", err)
	}
	return nil
}

func (s *Service) EditUserInTeam(username string, teamId *int64) error {
	userE, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	err = s.repo.UpdateUserInTeam(userE.Id, teamId)
	if err != nil {
		return fmt.Errorf("error edditing user to team: %w", err)
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

func (s *Service) GetTeamForUsername(ctx context.Context) (*model.ResponseTeamDTO, error) {
	user, err := s.getUserFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting user from ctx: %w", err)
	}
	if user.TeamId == nil {
		return nil, nil
	}
	entity, err := s.repo.GetTeamById(*user.TeamId)
	if err != nil {
		return nil, fmt.Errorf("error getting team: %w", err)
	}
	dto := infrastructure.MapTeamEntityToResponseDto(entity)
	return &dto, nil
}

func (s *Service) getUserFromCtx(ctx context.Context) (*model.UserEntity, error) {
	claims, ok := ctx.Value("user").(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("error getting claims")
	}
	user, err := s.repo.GetUserByUsername(claims.Username)
	if err != nil {
		return nil, fmt.Errorf("error getting user by username %w", err)
	}
	return &user, nil
}
