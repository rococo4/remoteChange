package model

type TeamCreateDto struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
}

type TeamEntity struct {
	Id          int64  `db:"id"`
	Name        string `db:"Name"`
	ClusterName string `db:"ClusterName"`
	Namespace   string `db:"Namespace"`
}

type UpdateTeamDTO struct {
	Id          int64   `json:"id"`
	Name        *string `json:"name"`
	ClusterName *string `json:"cluster_name"`
	Namespace   *string `json:"namespace"`
}
type ResponseTeamDTO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
}
