package team

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fridrock/projects_service/api"
)

type UsersClient interface {
	GetProfilesByIds(api.ProfilesByUserIdsDto) ([]api.GetUserDto, error)
}

type UsersClientImpl struct {
	remoteAddress string
}

func (uc *UsersClientImpl) GetProfilesByIds(dto api.ProfilesByUserIdsDto) ([]api.GetUserDto, error) {
	var profiles []api.GetUserDto
	param := "ids="
	for i := 0; i < len(dto.Ids); i++ {
		param += dto.Ids[i].String()
		if i != len(dto.Ids)-1 {
			param += ","
		}
	}
	fullURL := uc.remoteAddress + "/users/profiles?" + param
	response, err := http.Get(fullURL)
	if err != nil {
		return profiles, err
	}
	defer response.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return profiles, err
	}

	// Парсинг JSON-ответа в список DTO
	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return profiles, err
	}
	return profiles, nil
}

func NewUsersClient() UsersClient {
	varName := "USERS_ADDRESS"
	address, exists := os.LookupEnv(varName)
	if !exists {
		log.Fatalf("Can't load env variable: %v", varName)
	}
	return &UsersClientImpl{
		remoteAddress: address,
	}
}
