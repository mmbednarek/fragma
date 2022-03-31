package repo

import (
	coreV1Detail "github.com/mmbednarek/fragma/api/fragma/core/v1/detail"
)

func GetStandardRepository() ApiRepository {
	repo := ApiRepository{
		Apis: map[string]ApiDetail{},
	}
	repo.RegisterApi(coreV1Detail.ApiDetail{})
	return repo
}
