package infrastructure

import "remoteChange/internal/model"

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
		Name:        dto.Name,
		ClusterName: dto.ClusterName,
		Namespace:   dto.Namespace,
	}
}

func MapTeamEntityToResponseDto(entity model.TeamEntity) model.ResponseTeamDTO {
	return model.ResponseTeamDTO{
		Id:          entity.Id,
		Name:        entity.Name,
		ClusterName: entity.ClusterName,
		Namespace:   entity.Namespace,
	}
}

func MapCreateConfigDtoToConfigEntity(dto model.CreateConfigDTO) model.ConfigEntity {
	return model.ConfigEntity{
		TeamId:    dto.TeamId,
		Content:   dto.Content,
		CreatedAt: dto.CreatedAt,
	}
}

func MapConfigEntityToConfigResponse(entity model.ConfigEntity) model.ConfigResponse {
	return model.ConfigResponse{
		Id:        entity.Id,
		Name:      entity.Name,
		TeamId:    entity.TeamId,
		Type:      entity.Type,
		Content:   entity.Content,
		CreatedAt: entity.CreatedAt,
	}
}
