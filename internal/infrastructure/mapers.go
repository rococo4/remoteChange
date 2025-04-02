package infrastructure

import (
	"remoteChange/internal/model"
	"time"
)

func MapUserDtoRegisterToUserEntity(dto model.UserDTORegister) model.UserEntity {
	return model.UserEntity{
		Username: dto.Username,
		Name:     dto.Name,
		Surname:  dto.Surname,
		TeamId:   nil,
		Role:     "user",
		Password: dto.Password,
	}
}

func MapTeamDtoCreateToTeamEntity(dto model.TeamCreateDto) model.TeamEntity {
	return model.TeamEntity{
		Name: dto.Name,
	}
}

func MapTeamEntityToResponseDto(entity model.TeamEntity) model.ResponseTeamDTO {
	return model.ResponseTeamDTO{
		Id:   entity.Id,
		Name: entity.Name,
	}
}

func MapCreateConfigDtoToConfigEntity(dto model.CreateConfigDTO) model.ConfigEntity {
	return model.ConfigEntity{
		TeamId:      dto.TeamId,
		Content:     dto.Content,
		CreatedAt:   time.Now(),
		Description: dto.Description,
	}
}

func MapConfigEntityToConfigResponse(entity model.ConfigEntity) model.ConfigResponse {
	return model.ConfigResponse{
		Id:          entity.Id,
		Name:        entity.Name,
		TeamId:      entity.TeamId,
		Type:        entity.Type,
		Content:     entity.Content,
		CreatedAt:   entity.CreatedAt,
		Description: entity.Description,
	}
}

type repo interface {
	GetUserById(id int64) (model.UserEntity, error)
}

func MapConfigChangesEntityToConfigVersionResponse(entity model.ConfigChangesEntity, r repo) (model.ConfigVersionResponse, error) {
	user, err := r.GetUserById(entity.UserId)
	if err != nil {
		return model.ConfigVersionResponse{}, err
	}
	return model.ConfigVersionResponse{
		Id:       entity.Id,
		Username: user.Username,
		Action:   entity.Action,
		Date:     entity.ActionAt,
	}, nil
}
