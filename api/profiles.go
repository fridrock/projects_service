package api

import "github.com/google/uuid"

type ProfilesByUserIdsDto struct {
	Ids []uuid.UUID `json:"ids"`
}
